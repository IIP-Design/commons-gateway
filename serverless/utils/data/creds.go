package data

import (
	"fmt"
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
		logError(err)
	}

	// creds := map[string]interface{}{
	// 	"hash": pass_hash,
	// 	"salt": salt,
	// }

	creds := CredentialsData{
		Hash: pass_hash,
		Salt: salt,
	}

	return creds, err
}
