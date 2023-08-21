package main

import (
	"context"
	"errors"

	"github.com/IIP-Design/commons-gateway/utils/data/admins"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
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

	err = admins.CreateAdmin(adminData)

	return err
}

// NewAdminHandler handles the request to create a new administrative user. It
// ensures that the required data is present before continuing on to recording
// the user's email in the list of admins.
func NewAdminHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	var err error

	admin, err := data.ExtractUser(event.Body)

	if err != nil {
		return msgs.SendServerError(err)
	}

	err = handleAdminCreation(admin)

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.PrepareResponse([]byte("success"))
}

func main() {
	lambda.Start(NewAdminHandler)
}
