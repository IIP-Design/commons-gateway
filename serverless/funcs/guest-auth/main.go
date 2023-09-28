package main

import (
	"context"
	"errors"
	"os"

	"github.com/IIP-Design/commons-gateway/utils/data/creds"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
	"github.com/IIP-Design/commons-gateway/utils/security/jwt"
	"github.com/IIP-Design/commons-gateway/utils/turnstile"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// verify2FA retrieves that the user provided 2FA
// code matches the values sent to the user.
func verify2FA(id string, code string) bool {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	var storedCode string

	query := "SELECT code FROM mfa WHERE request_id = $1;"
	err = pool.QueryRow(query, id).Scan(&storedCode)

	if err != nil {
		logs.LogError(err, "Retrieve MFA Query Error")
		return false
	}

	return code == storedCode
}

// clear2FA removes the record of a 2FA request after it has been used.
func clear2FA(id string) {
	pool := data.ConnectToDB()
	defer pool.Close()

	query := "DELETE FROM mfa WHERE request_id = $1;"
	_, err := pool.Exec(query, id)

	if err != nil {
		logs.LogError(err, "Delete MFA Request Query Error")
	}
}

// handleGrantAccess ensures that a user hash provided a password has matching their
// username and if so, generates a JWT to grant them guest access.
func handleGrantAccess(username string, clientHash string) (msgs.Response, error) {
	if clientHash == "" || username == "" {
		return msgs.Response{StatusCode: 400}, errors.New("data missing from request")
	}

	creds, err := creds.RetrieveCredentials(username)

	if err != nil {
		return msgs.SendServerError(err)
	}

	if creds.Hash != clientHash {
		return msgs.SendAuthError(errors.New("forbidden"), 403)
	} else if creds.Expired {
		return msgs.SendAuthError(errors.New("credentials expired"), 403)
	} else if !creds.Approved {
		return msgs.SendAuthError(errors.New("user is not yet approved"), 403)
	}

	jwt, err := jwt.FormatJWT(username, creds.Role)

	if err != nil {
		return msgs.SendServerError(err)
	}

	body, err := msgs.MarshalBody(jwt)

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.PrepareResponse(body)
}

// authenticationHandler manages guest user authentication by either generating a JSON web
// token for new authentication sessions or verifying an existing token for an ongoing session.
func authenticationHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	// Retrieve the auth data submitted by the user.
	parsed, err := data.ParseBodyData(event.Body)

	if err != nil {
		return msgs.SendServerError(err)
	}

	clientHash := parsed.Hash
	username := parsed.Username
	mfaId := parsed.MFA.Id
	mfaCode := parsed.MFA.Code

	// Verify that the provided 2FA code is valid.
	verified := verify2FA(mfaId, mfaCode)

	if !verified {
		return msgs.SendAuthError(errors.New("forbidden"), 403)
	}

	// Delete the existing 2FA entry
	clear2FA(mfaId)

	// Verify the turnstile captcha token
	tokenVerSecretKey := os.Getenv("TOKEN_VERIFICATION_SECRET_KEY")

	if tokenVerSecretKey != "" {
		token := parsed.Token
		remoteIp := event.RequestContext.Identity.SourceIP

		valid, err := turnstile.TokenIsValid(token, remoteIp, tokenVerSecretKey)

		if !valid || err != nil {
			logs.LogError(err, "Turnstile error")
			return msgs.SendServerError(err)
		}
	}

	return handleGrantAccess(username, clientHash)
}

func main() {
	lambda.Start(authenticationHandler)
}
