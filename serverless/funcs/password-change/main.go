package main

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/IIP-Design/commons-gateway/utils/data/creds"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
	"github.com/IIP-Design/commons-gateway/utils/security/hashing"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nbutton23/zxcvbn-go"
)

type PasswordReset struct {
	CurrentPasswordHash string `json:"currentPasswordHash"`
	NewPassword         string `json:"newPassword"`
	Email               string `json:"email"`
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

func verifyUser(parsed PasswordReset) (creds.CredentialsData, error) {
	var credentials creds.CredentialsData

	_, exists, err := data.CheckForExistingUser(parsed.Email, "guests")
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

func checkNewPassword(email string, role string, newPassword string) bool {
	result := zxcvbn.PasswordStrength(newPassword, []string{email, role})
	success := (result.Score >= 3)
	return success
}

func updatePassword(email string, salt string, newPassword string) error {
	var err error
	passHash := hashing.GenerateHash(newPassword, salt)

	pool := data.ConnectToDB()
	defer pool.Close()

	query := "UPDATE invites SET pass_hash = $1, first_login = FALSE " +
		" WHERE invitee = $2 AND salt = $3 " +
		" AND pending = FALSE AND expiration > NOW() " +
		" AND date_invited = ( SELECT max(date_invited) FROM invites WHERE invitee = $2 AND pending = FALSE )"

	_, err = pool.Exec(query, passHash, email, salt)

	return err
}

func PasswordChangeHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	parsed, err := extractBody(event.Body)
	if err != nil {
		return msgs.SendServerError(err)
	}

	credentials, err := verifyUser(parsed)
	if err != nil {
		return msgs.SendServerError(err)
	}

	passwordIsSecure := checkNewPassword(parsed.Email, credentials.Role, parsed.NewPassword)
	if !passwordIsSecure {
		return msgs.SendCustomError(errors.New("password is not strong enough"), 422)
	}

	err = updatePassword(parsed.Email, credentials.Salt, parsed.NewPassword)
	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(PasswordChangeHandler)
}
