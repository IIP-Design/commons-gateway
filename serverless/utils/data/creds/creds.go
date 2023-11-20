package creds

import (
	"errors"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/invites"
	"github.com/IIP-Design/commons-gateway/utils/data/users"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/IIP-Design/commons-gateway/utils/security/hashing"
)

type CredentialsData struct {
	Hash       string   `json:"hash"`
	Salt       string   `json:"salt"`
	PrevSalts  []string `json:"prevSalts"`
	Expired    bool     `json:"expired"`
	Approved   bool     `json:"approved"`
	Locked     bool     `json:"locked"`
	FirstLogin bool     `json:"firstLogin"`
	Role       string   `json:"role"`
}

// ClearUnsuccessfulLoginAttempts resets the given user's login counter to zero.
func ClearUnsuccessfulLoginAttempts(guest string) error {
	pool := data.ConnectToDB()
	defer pool.Close()

	query := `UPDATE guests SET login_attempt = 0, login_date = NULL WHERE email = $1;`
	_, err := pool.Exec(query, guest)

	if err != nil {
		logs.LogError(err, "Clear Login Attempts Query Error")
	}

	return err
}

// RetrieveCredentials
func RetrieveCredentials(email string) (CredentialsData, error) {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	var passHash string
	var salt string
	var prevSalts []string
	var expired bool
	var approved bool
	var locked bool
	var firstLogin bool
	var role string

	query :=
		`SELECT pass_hash, salt, expiration < NOW() AS expired, pending=FALSE AS approved, locked, first_login, role
		 FROM guest_auth_data WHERE email = $1;`

	err = pool.QueryRow(query, email).Scan(&passHash, &salt, &expired, &approved, &locked, &firstLogin, &role)

	if err != nil {
		logs.LogError(err, "Retrieve Credentials Query Error")
	}

	rows, err := pool.Query(`SELECT salt FROM password_history WHERE user_id = $1;`, email)

	if err != nil {
		logs.LogError(err, "Get Previous Salts Query Error")
	}

	defer rows.Close()

	for rows.Next() {
		var salt string

		if err := rows.Scan(&salt); err != nil {
			logs.LogError(err, "Get Salt Query Error")
		}

		prevSalts = append(prevSalts, salt)
	}

	creds := CredentialsData{
		Hash:       passHash,
		Salt:       salt,
		PrevSalts:  prevSalts,
		Expired:    expired,
		Approved:   approved,
		Locked:     locked,
		FirstLogin: firstLogin,
		Role:       role,
	}

	return creds, err
}

func SaveInitialInvite(invite data.Invite, setPending bool) (string, error) {
	var pass string

	// Ensure invitee doesn't already have access.
	_, guestHasAccess, err := users.CheckForExistingUser(invite.Invitee.Email, "guests")

	if err != nil {
		return pass, err
	} else if guestHasAccess {
		return pass, errors.New("guest user already exists")
	}

	// Save credentials
	err = invites.SaveCredentials(invite.Invitee)

	if err != nil {
		return pass, errors.New("something went wrong - credential generation failed")
	}

	// PASSWORD IS UNRECOVERABLE
	pass, salt := hashing.GenerateCredentials()
	hash := hashing.GenerateHash(pass, salt)

	// Record the invitation - has to follow cred generation due to foreign key constraint
	var email string

	if setPending {
		email = invite.Proposer
	} else {
		email = invite.Inviter
	}

	err = invites.SaveInvite(email, invite.Invitee.Email, invite.Expires, hash, salt, setPending, true, true)

	if err != nil {
		return pass, errors.New("something went wrong - saving invite failed")
	}

	return pass, nil
}

func ResetPassword(email string) (string, error) {
	pool := data.ConnectToDB()
	defer pool.Close()

	pass, salt := hashing.GenerateCredentials()
	hash := hashing.GenerateHash(pass, salt)

	query :=
		`UPDATE invites SET salt = $1, pass_hash = $2 WHERE invitee = $3 AND date_invited = ( SELECT MAX(date_invited)
		 FROM invites WHERE invitee = $3 AND pending = FALSE );`
	_, err := pool.Exec(query, salt, hash, email)

	if err != nil {
		logs.LogError(err, "Reset Password Error")
	}

	return pass, err
}
