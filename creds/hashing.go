package main

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/pbkdf2"
)

func GeneratePassword() (string, string) {
	pass := randstr.String(20)
	salt := randstr.String(10)

	return pass, salt
}

func GenerateHash(pass string, salt string) string {
	var iterations = 4096

	derivedKey := pbkdf2.Key([]byte(pass), []byte(salt), iterations, 32, sha256.New)

	encoded := base64.StdEncoding.EncodeToString(derivedKey)

	return encoded
}
