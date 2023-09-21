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

// AuthenticationHandler manages guest user authentication by either generating a JSON web
// token for new authentication sessions or verifying an existing token for an ongoing session.
func AuthenticationHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	parsed, err := data.ParseBodyData(event.Body)

	if err != nil {
		return msgs.SendServerError(err)
	}

	action := parsed.Action
	clientHash := parsed.Hash
	username := parsed.Username

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

	if action == "create" {
		return handleGrantAccess(username, clientHash)
	}

	return msgs.Response{StatusCode: 400}, errors.New("invalid request")
}

func main() {
	lambda.Start(AuthenticationHandler)
}
