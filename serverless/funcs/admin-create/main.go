package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/IIP-Design/commons-gateway/utils/data/admins"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/users"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
	"github.com/IIP-Design/commons-gateway/utils/types"
)

// handleAdminCreation coordinates all the actions associated with creating a new user.
func handleAdminCreation(adminData types.User) (bool, error) {
	var err error

	_, isAdmin, err := users.CheckForExistingUser(adminData.Email, "admins")

	if err != nil {
		return isAdmin, err
	} else if isAdmin {
		return isAdmin, errors.New("this user has already been added as an administrator")
	}

	err = admins.CreateAdmin(adminData)

	return isAdmin, err
}

// newAdminHandler handles the request to create a new administrative user. It
// ensures that the required data is present before continuing on to recording
// the user's email in the list of admins.
func newAdminHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	admin, err := data.ExtractUser(event.Body)

	if err != nil {
		return msgs.SendServerError(err)
	}

	exists, err := handleAdminCreation(admin)

	if exists {
		return msgs.SendCustomError(err, 409)
	} else if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(newAdminHandler)
}
