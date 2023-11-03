package hashing

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/IIP-Design/commons-gateway/utils/randstr"
	"golang.org/x/crypto/pbkdf2"
)

const (
	PASS_LEN = 20
	SALT_LEN = 10
)

// generateCredentials creates a random PASS_LEN character string to be used as a password
// as well as a random SALT_LEN character string to salt the password when hashing
// it for storage in the database.
func GenerateCredentials() (string, string) {
	pass, err := randstr.RandStringBytes(PASS_LEN)

	if err != nil {
		logs.LogError(err, "failed to generate credentials - password")
	}

	salt, err := randstr.RandStringBytes(SALT_LEN)

	if err != nil {
		logs.LogError(err, "failed to generate credentials - salt")
	}

	return pass, salt
}

// generateHash returns a base64-encoded hash of the provided password and salt values.
// The salt is appended to the password and the combination is run through 4096 iterations
// of PBKDF2 using the SHA-256 hashing function. The resulting 32 byte derived key is then
// encoded as a base64 string for ease of use.
func GenerateHash(pass string, salt string) string {
	var iterations = 4096
	var keyLength = 32

	derivedKey := pbkdf2.Key([]byte(pass), []byte(salt), iterations, keyLength, sha256.New)

	encoded := base64.StdEncoding.EncodeToString(derivedKey)

	return encoded
}
