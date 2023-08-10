package data

import (
	"fmt"

	"github.com/IIP-Design/commons-gateway/utils/logs"
)

type CredentialsData struct {
	Hash string `json:"hash"`
	Salt string `json:"salt"`
}

// RetrieveCredentials
func RetrieveCredentials(email string) (CredentialsData, error) {
	var err error

	pool := connectToDB()
	defer pool.Close()

	var pass_hash string
	var salt string

	query := fmt.Sprintf(`SELECT pass_hash, salt FROM credentials WHERE email = '%s';`, email)
	err = pool.QueryRow(query).Scan(&pass_hash, &salt)

	if err != nil {
		logs.LogError(err, "Retrieve Credentials Query Error")
	}

	creds := CredentialsData{
		Hash: pass_hash,
		Salt: salt,
	}

	return creds, err
}
