package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/IIP-Design/commons-gateway/utils/data/creds"
	"github.com/IIP-Design/commons-gateway/utils/data/users"
	"github.com/IIP-Design/commons-gateway/utils/email/provision"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
)

// PasswordResetHandler handles the request to retrieve a single admin user based on email address.
func PasswordResetHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	id := event.QueryStringParameters["id"]

	if id == "" {
		return msgs.SendServerError(errors.New("user id not provided"))
	}

	// Ensure the user exists doesn't already have access.
	pass, err := creds.ResetPassword(id)

	if err != nil {
		return msgs.SendServerError(err)
	}

	// Ensure the user exists doesn't already have access.
	user, _, err := users.CheckForExistingUser(id, "guests")

	if err != nil {
		return msgs.SendServerError(err)
	}

	_, err = provision.MailProvisionedCreds(user, pass, 3)
	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(PasswordResetHandler)
}
