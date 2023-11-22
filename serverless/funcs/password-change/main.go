package main

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/xid"

	"github.com/IIP-Design/commons-gateway/utils/data/creds"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/users"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
)

type PasswordReset struct {
	CurrentPasswordHash string   `json:"currentPasswordHash"`
	Email               string   `json:"email"`
	HashedPriorSalts    []string `json:"hashesWithPriorSalts"`
	NewPasswordHash     string   `json:"newPasswordHash"`
	NewSalt             string   `json:"newSalt"`
}

func extractBody(body string) (PasswordReset, error) {
	var parsed PasswordReset

	b := []byte(body)
	err := json.Unmarshal(b, &parsed)

	if err != nil {
		logs.LogError(err, "Failed to Unmarshal Body")
	}

	return parsed, err
}

// verifyUser confirms that the user requesting a password change exists
// and has provided the correct password.
func verifyUser(parsed PasswordReset) (creds.CredentialsData, error) {
	var credentials creds.CredentialsData

	_, exists, err := users.CheckForExistingUser(parsed.Email, "guests")

	if err != nil || !exists {
		return credentials, errors.New("user does not exist")
	}

	credentials, err = creds.RetrieveCredentials(parsed.Email)

	if err != nil {
		return credentials, errors.New("failed to load credentials")
	}

	if credentials.Hash != parsed.CurrentPasswordHash {
		return credentials, errors.New("credentials do not match")
	}

	return credentials, nil
}

// checkPasswordReused compares a list of provided password hashes (generally a new
// password hashed with the salts from previous passwords) against a list of a user's
// previous password hashes. A match between these lists indicates password reuse.
func checkPasswordReused(email string, hashedPriorSalts []string) (bool, error) {
	var err error
	reused := false

	pool := data.ConnectToDB()
	defer pool.Close()

	query := "SELECT creation_date, pass_hash FROM password_history WHERE user_id = $1 ORDER BY creation_date DESC LIMIT 24;"
	rows, err := pool.Query(query, email)

	if err != nil {
		logs.LogError(err, "Get Guests Query Error")
		return reused, err
	}

	// Convert list of hashes generated using previous salts into a map for easier searching.
	hashMap := make(map[string]string, len(hashedPriorSalts))
	for i := range hashedPriorSalts {
		hashMap[hashedPriorSalts[i]] = hashedPriorSalts[i]
	}

	defer rows.Close()

	for rows.Next() {
		var date time.Time
		var passHash string

		if err := rows.Scan(&date, &passHash); err != nil {
			logs.LogError(err, "Pass Reuse Scan Error")
			return reused, err
		}

		_, found := hashMap[passHash]

		if found {
			reused = true
			logs.LogError(errors.New("matching hash found"), "Password Hash Collision Error")
			break
		}
	}

	if err = rows.Err(); err != nil {
		logs.LogError(err, "Pass Reuse Scan Error")
	}

	return reused, err
}

// updatePassword stores a hash of the newly created password as needed in the database.
func updatePassword(email string, salt string, newPasswordHash string, newSalt string) error {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	// Save new credentials to invites table.
	query := "UPDATE invites SET pass_hash = $1, salt = $2, first_login = FALSE " +
		" WHERE invitee = $3 AND salt = $4 " +
		" AND pending = FALSE AND expiration > NOW() " +
		" AND date_invited = ( SELECT max(date_invited) FROM invites WHERE invitee = $3 AND pending = FALSE )"
	_, err = pool.Exec(query, newPasswordHash, newSalt, email, salt)

	if err != nil {
		return err
	}

	// Save new credentials to password history table.
	id := xid.New()
	query = "INSERT INTO password_history ( id, user_id, creation_date, salt, pass_hash ) VALUES ( $1, $2, NOW(), $3, $4)"
	_, err = pool.Exec(query, id, email, newSalt, newPasswordHash)

	if err != nil {
		return err
	}

	// Limits the number of history entries stored per user.
	query = "DELETE FROM password_history WHERE user_id = $1 AND id NOT IN ( SELECT id FROM password_history WHERE user_id = $1 ORDER BY creation_date DESC LIMIT 24 )"
	_, err = pool.Exec(query, email)

	if err != nil {
		return err
	}

	return err
}

func passwordChangeHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	parsed, err := extractBody(event.Body)

	if err != nil {
		return msgs.SendServerError(err)
	}

	credentials, err := verifyUser(parsed)

	if err != nil {
		return msgs.SendServerError(err)
	}

	passwordIsReused, err := checkPasswordReused(parsed.Email, parsed.HashedPriorSalts)

	if err != nil {
		return msgs.SendServerError(err)
	} else if passwordIsReused {
		return msgs.SendCustomError(errors.New("password was reused"), 409)
	}

	err = updatePassword(parsed.Email, credentials.Salt, parsed.NewPasswordHash, parsed.NewSalt)

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(passwordChangeHandler)
}
