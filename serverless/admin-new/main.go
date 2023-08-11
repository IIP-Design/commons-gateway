package main

import (
	"context"
	"errors"

	data "github.com/IIP-Design/commons-gateway/utils/data"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// handleAdminCreation coordinates all the actions associated with creating a new user.
func handleAdminCreation(adminEmail string) error {
	var err error

	isAdmin, err := data.CheckForExistingUser(adminEmail, "admin")

	if err != nil {
		return err
	} else if isAdmin {
		return errors.New("this user has already been added as an administrator")
	}

	err = data.CreateAdmin(adminEmail)

	return err
}

// NewAdminHandler handles the request to create a new administrative user. It
// ensures that the required data is present before continuing on to recording
// the user's email in the list of admins.
func NewAdminHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	var msg string

	parsed, err := data.ParseBodyData(event.Body)

	adminEmail := parsed.Email

	if err != nil {
		return msgs.Response{StatusCode: 500}, err
	} else if adminEmail == "" {
		return msgs.Response{StatusCode: 400}, errors.New("data missing from request")
	}

	err = handleAdminCreation(adminEmail)

	if err != nil {
		return msgs.Response{StatusCode: 500}, err
	} else {
		msg = "success"
	}

	return msgs.PrepareResponse(msg)
}

func main() {
	lambda.Start(NewAdminHandler)
}
