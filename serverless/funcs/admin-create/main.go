package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/IIP-Design/commons-gateway/utils/data/admins"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/users"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
)

// handleAdminCreation coordinates all the actions associated with creating a new user.
func handleAdminCreation(adminData data.User) (bool, error) {
	var err error

	exists, user, err := users.CheckForExistingUser(adminData.Email)

	if err != nil {
		logs.LogError(err, "Check For Existing User Error")
		return exists, err
	} else if exists {
		err = fmt.Errorf("the user %s has already been registered as a user of type %s", adminData.Email, user.Type)

		logs.LogError(err, "Check For Existing User Error")
		return exists, err
	}

	err = admins.CreateAdmin(adminData)

	if err != nil {
		logs.LogError(err, "Admin Creation Error")
	}

	return exists, err
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
