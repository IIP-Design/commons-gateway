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
func handleAdminCreation(adminData data.User) error {
	var err error

	isAdmin, err := data.CheckForExistingUser(adminData.Email, "admins")

	if err != nil {
		return err
	} else if isAdmin {
		return errors.New("this user has already been added as an administrator")
	}

	err = data.CreateAdmin(adminData)

	return err
}

// NewAdminHandler handles the request to create a new administrative user. It
// ensures that the required data is present before continuing on to recording
// the user's email in the list of admins.
func NewAdminHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	var msg string
	var err error

	admin, err := data.ExtractUser(event.Body)

	if err != nil {
		return msgs.Response{StatusCode: 500}, err
	}

	err = handleAdminCreation(admin)

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
