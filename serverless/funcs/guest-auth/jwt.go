package main

import (
	"encoding/json"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var secret = []byte(os.Getenv("JWT_SECRET"))

// generateJWT creates a JSON web token that can be used to authenticate to the
// web application. The token contains the user's name and access scope and is
// valid for one hour.
func generateJWT(username string, scope string) (string, error) {
	// TODO: Switch to EdDSA signing key
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["exp"] = time.Now().Add(1 * time.Hour)
	claims["scope"] = scope
	claims["user"] = username

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	webToken, err := json.Marshal(
		map[string]interface{}{
			"token": tokenString,
		},
	)

	if err != nil {
		return "", err
	}

	return string(webToken), nil
}
