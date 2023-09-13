package creds

import (
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
)

type CredentialsData struct {
	Hash     string `json:"hash"`
	Salt     string `json:"salt"`
	Expired  bool   `json:"expired"`
	Approved bool   `json:"approved"`
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

	query := `SELECT pass_hash, salt, expiration < NOW() AS expired, pending=FALSE AS approved FROM guests LEFT JOIN invites ON guests.email=invites.invitee WHERE email = $1;`
	err = pool.QueryRow(query, email).Scan(&pass_hash, &salt, &expired, &approved)

	if err != nil {
		logs.LogError(err, "Retrieve Credentials Query Error")
	}

	creds := CredentialsData{
		Hash:     pass_hash,
		Salt:     salt,
		Expired:  expired,
		Approved: approved,
	}

	return creds, err
}
