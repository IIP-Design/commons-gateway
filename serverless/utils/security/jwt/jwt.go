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

func VerifyJWT(tokenString string, scopes []string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if err != nil {
		logs.LogError(err, "Error Parsing JWT Token")
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		logs.LogError(err, "Bearer Token is Not Valid")
		return errors.New("token is not valid")
	}

	if !slices.Contains(scopes, claims["scope"].(string)) {
		logs.LogError(err, "Bearer Token Has Incorrect Scope")
		return errors.New("token has incorrect scope: " + claims["scope"].(string))
	}

	return nil
}

// DEPRECATED - remove once all endpoints are switched over to using authorizer function.
func RequestIsAuthorized(req events.APIGatewayProxyRequest, scopes []string) (int, error) {
	authHeader := req.Headers["Authorization"]
	token, err := extractBearerToken(authHeader)

	if err != nil {
		return 401, err
	}

	err = VerifyJWT(token, scopes)

	// Everything looks good and client can continue
	if err == nil {
		return 200, nil
		// Token is expired and client should re-login
	} else if errors.Is(err, jwt.ErrTokenExpired) {
		return 401, err
		// Token is invalid and client must be rejected
	} else {
		return 403, err
	}
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
