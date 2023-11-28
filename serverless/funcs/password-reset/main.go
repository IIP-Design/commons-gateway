package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/IIP-Design/commons-gateway/utils/data/creds"
	"github.com/IIP-Design/commons-gateway/utils/data/users"
	"github.com/IIP-Design/commons-gateway/utils/email/provision"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
)

// passwordResetHandler handles the request to retrieve a single admin user based on email address.
func passwordResetHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	id := event.QueryStringParameters["id"]

	if id == "" {
		return msgs.SendCustomError(errors.New("user id not provided"), 400)
	}

	// Ensure the user exists
	user, exists, err := users.CheckForExistingGuestUser(id)

	if err != nil {
		logs.LogError(err, "Check For Guest User Error")
		return msgs.SendServerError(err)
	} else if !exists {
		err = fmt.Errorf("%s is not registered as a guest user", id)

		logs.LogError(err, "Guest User Not Found Error")
		return msgs.SendCustomError(errors.New("user does not exist"), 404)
	}

	pass, err := creds.ResetPassword(id)

	if err != nil {
		logs.LogError(err, "Reset Password Error")
		return msgs.SendServerError(err)
	}

	_, err = provision.MailProvisionedCreds(user, pass, 2)

	if err != nil {
		logs.LogError(err, "Mail Credentials Error")
		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(passwordResetHandler)
}
