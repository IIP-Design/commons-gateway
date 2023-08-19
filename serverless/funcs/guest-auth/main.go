package main

import (
	"errors"

	"github.com/IIP-Design/commons-gateway/utils/data/creds"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

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
func AuthenticationHandler(ctx events.APIGatewayProxyRequestContext, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	parsed, err := data.ParseBodyData(event.Body)

	if err != nil {
		return msgs.SendServerError(err)
	}

	action := parsed.Action
	clientHash := parsed.Hash
	username := parsed.Username

	if action == "create" {
		return handleGrantAccess(username, clientHash)
	}

	return msgs.Response{StatusCode: 400}, errors.New("invalid request")
}

func main() {
	lambda.Start(AuthenticationHandler)
}
