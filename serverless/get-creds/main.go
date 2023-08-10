package main

import (
	"context"
	"encoding/json"
	"errors"

	data "github.com/IIP-Design/commons-gateway/utils/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/lambda"
)

// EventData describes the data that the Lambda function expects to receive.
type EventData struct {
	Username string `json:"username"`
}

// handleCredentialRequest coordinates all the actions associated with retrieving user credentials.
func handleCredentialRequest(username string) (data.CredentialsData, error) {
	var err error
	var creds data.CredentialsData

	exists, err := data.CheckForExistingUser(username, "credentials")

	if err != nil {
		return creds, err
	} else if !exists {
		return creds, errors.New("user not found")
	}

	creds, err = data.RetrieveCredentials(username)

	return creds, err
}

// GetCredsHandler handles the request to retrieve the password hash and salt associated
// with a user based on the user name.
func GetCredsHandler(ctx context.Context, event EventData) (msgs.Response, error) {
	var msg []byte

	user := event.Username

	if user == "" {
		logs.LogError(nil, "Username not provided in request.")
		return msgs.Response{StatusCode: 400}, errors.New("data missing from request")
	}

	creds, err := handleCredentialRequest(user)

	if err != nil {
		return msgs.Response{StatusCode: 500}, err
	} else {
		msg, err = json.Marshal(creds)

		if err != nil {
			return msgs.Response{StatusCode: 500}, err
		}
	}

	return msgs.PrepareResponse(string(msg))
}

func main() {
	lambda.Start(GetCredsHandler)
}
