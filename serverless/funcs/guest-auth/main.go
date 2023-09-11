package main

import (
	"context"
	"errors"
	"os"

	"github.com/IIP-Design/commons-gateway/utils/data/creds"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
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
		return msgs.Response{Body: "Forbidden", StatusCode: 403}, err
	}

	jwt, err := generateJWT(username, "guest")

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
		err := turnstile.VerifyToken(token, remoteIp, tokenVerSecretKey)
		if err != nil {
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
