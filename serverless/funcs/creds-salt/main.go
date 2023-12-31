package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/IIP-Design/commons-gateway/utils/data/creds"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/users"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// handleCredentialRequest coordinates all the actions associated with retrieving user credentials.
func handleCredentialRequest(username string) (creds.CredentialsData, error) {
	var err error
	var credentials creds.CredentialsData

	_, exists, err := users.CheckForExistingGuestUser(username)

	if err != nil {
		logs.LogError(err, "Check For Guest User Error")
		return credentials, err
	} else if !exists {
		err = fmt.Errorf("%s is not registered as a guest user", username)

		logs.LogError(err, "Guest User Not Found Error")
		return credentials, errors.New("user not found")
	}

	credentials, err = creds.RetrieveCredentials(username)

	return credentials, err
}

// getSaltHandler handles the request to retrieve the salt associated with a user based on the user name.
func getSaltHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	parsed, err := data.ParseBodyData(event.Body)

	user := parsed.Username

	if err != nil {
		return msgs.SendServerError(err)
	} else if user == "" {
		err = errors.New("data missing from request")
		logs.LogError(err, "Username not provided in request.")
		return msgs.SendCustomError(err, 400)
	}

	credentials, err := handleCredentialRequest(user)

	if err != nil {
		return msgs.SendServerError(err)
	}

	if credentials.Locked {
		err = errors.New("account locked")
		logs.LogError(err, "User's account is locked.")
		return msgs.SendCustomError(err, 429)
	}

	salts := map[string]any{
		"salt":      credentials.Salt,
		"prevSalts": credentials.PrevSalts,
	}

	body, err := msgs.MarshalBody(salts)

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.PrepareResponse(body)
}

func main() {
	lambda.Start(getSaltHandler)
}
