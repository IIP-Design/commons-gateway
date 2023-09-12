package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	jwt "github.com/golang-jwt/jwt/v5"
)

var secret = []byte(os.Getenv("JWT_SECRET"))

// generateJWT creates a JSON web token that can be used to authenticate to the
// web application. The token contains the user's name and access scope and is
// valid for one hour.
func GenerateJWT(username string, scope string) (string, error) {
	// TODO: Switch to EdDSA signing key
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["exp"] = time.Now().Add(1 * time.Hour).Unix()
	claims["scope"] = scope
	claims["user"] = username

	return token.SignedString(secret)
}

func FormatJWT(username string, scope string) (string, error) {
	tokenString, err := GenerateJWT(username, scope)

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

func ExtractBearerToken(headerVal string) (string, error) {
	if headerVal == "" {
		return "", errors.New("no bearer token received")
	}

	segments := strings.Split(headerVal, " ")
	if len(segments) == 1 {
		return segments[0], nil
	} else if len(segments) == 2 {
		return segments[1], nil
	} else {
		return headerVal, errors.New("unable to extract bearer token")
	}
}

func VerifyJWT(tokenString string, scopes []string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return errors.New("token is not valid")
	}

	if !slices.Contains(scopes, claims["scope"].(string)) {
		return errors.New("token has incorrect scope: " + claims["scope"].(string))
	}

	return nil
}

func RequestIsAuthorized(req events.APIGatewayProxyRequest, scopes []string) (bool, error) {
	authHeader := req.Headers["Authorization"]
	token, err := ExtractBearerToken(authHeader)
	if err != nil {
		return false, err
	}

	err = VerifyJWT(token, scopes)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
