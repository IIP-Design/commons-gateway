package main

import (
	"context"
	"errors"

	data "github.com/IIP-Design/commons-gateway/utils/data"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/lambda"
)

// EventData describes the data that the Lambda function expects to receive.
type EventData struct {
	Email string `json:"email"`
}

// handleAdminCreation coordinates all the actions associated with creating a new user.
func handleAdminCreation(adminEmail string) error {
	var err error

	isExistingAdmin, err := data.CheckForExistingUser(adminEmail, "admins")

	if err != nil {
		return err
	}

	if isExistingAdmin {
		return errors.New("already an admin user")
	} else {
		// Record the invitation
		err = data.CreateAdmin(adminEmail)

		if err != nil {
			return errors.New("something went wrong - admin creation failed")
		}
	}

	return err
}

// ProvisionHandler handles the request to create a new administrative user. It
// ensures that the required data is present before continuing on to recording
// the user's email in the list of admins.
func ProvisionHandler(ctx context.Context, event EventData) (msgs.Response, error) {
	var msg string

	adminEmail := event.Email

	if adminEmail == "" {
		return msgs.Response{StatusCode: 400}, errors.New("data missing from request")
	}

	err := handleAdminCreation(adminEmail)

	if err != nil {
		return msgs.Response{StatusCode: 500}, err
	} else {
		msg = "success"
	}

	return msgs.PrepareResponse(msg)
}

func main() {
	lambda.Start(ProvisionHandler)
}
