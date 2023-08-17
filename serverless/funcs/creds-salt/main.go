package main

import (
	"context"
	"encoding/json"
	"errors"

	data "github.com/IIP-Design/commons-gateway/utils/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// handleCredentialRequest coordinates all the actions associated with retrieving user credentials.
func handleCredentialRequest(username string) (data.CredentialsData, error) {
	var err error
	var creds data.CredentialsData

	exists, err := data.CheckForExistingUser(username, "guests")

	if err != nil {
		return creds, err
	} else if !exists {
		return creds, errors.New("user not found")
	}

	creds, err = data.RetrieveCredentials(username)

	return creds, err
}

// GetSaltHandler handles the request to retrieve the salt associated with a user based on the user name.
func GetSaltHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	var msg []byte

	parsed, err := data.ParseBodyData(event.Body)

	user := parsed.Username

	if err != nil {
		return msgs.Response{StatusCode: 500}, err
	} else if user == "" {
		logs.LogError(nil, "Username not provided in request.")
		return msgs.Response{StatusCode: 400}, errors.New("data missing from request")
	}

	creds, err := handleCredentialRequest(user)

	if err != nil {
		return msgs.Response{StatusCode: 500}, err
	} else {
		msg, err = json.Marshal(creds.Salt)

		if err != nil {
			return msgs.Response{StatusCode: 500}, err
		}
	}

	return msgs.PrepareResponse(string(msg))
}

func main() {
	lambda.Start(GetSaltHandler)
}
