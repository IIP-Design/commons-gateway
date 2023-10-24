package creds

import (
	"errors"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/invites"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/IIP-Design/commons-gateway/utils/security/hashing"
)

type CredentialsData struct {
	Hash     string `json:"hash"`
	Salt     string `json:"salt"`
	Expired  bool   `json:"expired"`
	Approved bool   `json:"approved"`
	Locked   bool   `json:"locked"`
	Role     string `json:"role"`
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

	var pass_hash string
	var salt string
	var expired bool
	var approved bool
	var locked bool
	var role string

	query :=
		`SELECT pass_hash, salt, expiration < NOW() AS expired, pending=FALSE AS approved, locked, role
		 FROM guest_auth_data WHERE email = $1;`

	err = pool.QueryRow(query, email).Scan(&pass_hash, &salt, &expired, &approved, &locked, &role)

	if err != nil {
		logs.LogError(err, "Retrieve Credentials Query Error")
	}

	creds := CredentialsData{
		Hash:     pass_hash,
		Salt:     salt,
		Expired:  expired,
		Approved: approved,
		Locked:   locked,
		Role:     role,
	}

	return creds, err
}

func SaveInitialInvite(invite data.Invite, setPending bool) (string, error) {
	var pass string

	// Ensure invitee doesn't already have access.
	_, guestHasAccess, err := data.CheckForExistingUser(invite.Invitee.Email, "guests")

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
	err = invites.SaveInvite(invite.Proposer, invite.Invitee.Email, invite.Expires, hash, salt, setPending, true)

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

	query := `UPDATE invites SET salt = $1, pass_hash = $2 WHERE invitee = $3 AND date_invited = ( SELECT MAX(date_invited) FROM invites WHERE invitee = $3 AND pending = FALSE );`
	_, err := pool.Exec(query, salt, hash, email)

	if err != nil {
		logs.LogError(err, "Reset Password Error")
	}

	return pass, err
}
