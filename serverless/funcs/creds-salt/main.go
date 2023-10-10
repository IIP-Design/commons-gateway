package main

import (
	"context"
	"errors"

	"github.com/IIP-Design/commons-gateway/utils/data/creds"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// handleCredentialRequest coordinates all the actions associated with retrieving user credentials.
func handleCredentialRequest(username string) (creds.CredentialsData, error) {
	var err error
	var credentials creds.CredentialsData

	_, exists, err := data.CheckForExistingUser(username, "guests")

	if err != nil {
		return credentials, err
	} else if !exists {
		return credentials, errors.New("user not found")
	}

	credentials, err = creds.RetrieveCredentials(username)

	return credentials, err
}

// GetSaltHandler handles the request to retrieve the salt associated with a user based on the user name.
func GetSaltHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	parsed, err := data.ParseBodyData(event.Body)

	user := parsed.Username

	if err != nil {
		return msgs.SendServerError(err)
	} else if user == "" {
		logs.LogError(nil, "Username not provided in request.")
		return msgs.Response{StatusCode: 400}, errors.New("data missing from request")
	}

	credentials, err := handleCredentialRequest(user)

	if err != nil {
		return msgs.SendServerError(err)
	}

	if credentials.Locked {
		logs.LogError(nil, "User's account is locked.")
		return msgs.SendAuthError(errors.New("account locked"), 429)
	}

	body, err := msgs.MarshalBody(credentials.Salt)

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.PrepareResponse(body)
}

func main() {
	lambda.Start(GetSaltHandler)
}
