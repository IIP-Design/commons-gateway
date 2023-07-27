package main

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/pbkdf2"
)

// generateOTP creates a random 20 character string to be used as a password
// as well as a random 10 character string to salt the password when hashing
// it for storage in the database.
func generateOTP() (string, string) {
	pass := randstr.String(20)
	salt := randstr.String(10)

	return pass, salt
}

// generateHash returns a base64-encoded hash of the provided password and salt values.
// The salt is appended to the password and the combination is run through 4096 iterations
// of PBKDF2 using the SHA-256 hashing function. The resulting 32 byte derived key is then
// encoded as a base64 string for ease of use.
func generateHash(pass string, salt string) string {
	var iterations = 4096
	var keyLength = 32

	derivedKey := pbkdf2.Key([]byte(pass), []byte(salt), iterations, keyLength, sha256.New)

	encoded := base64.StdEncoding.EncodeToString(derivedKey)

	return encoded
}
