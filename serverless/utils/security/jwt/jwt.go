package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/IIP-Design/commons-gateway/utils/logs"
	jwt "github.com/golang-jwt/jwt/v5"
)

var secret = []byte(os.Getenv("JWT_SECRET"))

// generateJWT creates a JSON web token that can be used to authenticate to the
// web application. The token contains the user's name and access scope and is
// valid for one hour.
func GenerateJWT(username string, scope string, firstLogin bool) (string, error) {
	// TODO: Switch to EdDSA signing key
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["exp"] = time.Now().Add(1 * time.Hour).Unix()
	claims["scope"] = scope
	claims["user"] = username
	claims["firstLogin"] = firstLogin

	return token.SignedString(secret)
}

func FormatJWT(username string, scope string, firstLogin bool) (string, error) {
	tokenString, err := GenerateJWT(username, scope, firstLogin)

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

// extractBearerToken returns the token portion of an authorization header.
// Will function whether or not the token is preceded by the work `Bearer `.
func extractBearerToken(headerVal string) (string, error) {
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

func parseToken(tokenString string) (string, error) {
	var scope string

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if err != nil {
		logs.LogError(err, "Error Parsing JWT Token")
		return scope, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		logs.LogError(err, "Bearer Token is Not Valid")
		return scope, errors.New("token is not valid")
	}

	scope = claims["scope"].(string)

	return scope, err
}

func VerifyJWT(tokenString string, scopes []string) error {
	scope, err := parseToken(tokenString)
	if err != nil {
		return err
	}

	if !slices.Contains(scopes, scope) {
		logs.LogError(errors.New("scope error"), "Bearer Token Has Incorrect Scope")
		return errors.New("token has incorrect scope: " + scope)
	}

	return nil
}

// CheckAuthToken is used by the Authorizer function to extract the token
// in an API Gateway request's authorization header and then verify the
// validity of the extracted token.
func CheckAuthToken(token string, scopes []string) error {
	extracted, err := extractBearerToken(token)

	if err != nil {
		logs.LogError(err, "Error Extracting Bearer Token")
		return err
	}

	err = VerifyJWT(extracted, scopes)

	if err != nil {
		logs.LogError(err, "Error Verifying Bearer Token")
	}

	return err
}

func ExtractClientRole(token string) (string, error) {
	tokenString, err := extractBearerToken(token)
	if err != nil {
		logs.LogError(err, "Error Extracting Bearer Token")
		return "", err
	}

	return parseToken(tokenString)
}
