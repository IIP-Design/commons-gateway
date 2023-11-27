package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/IIP-Design/commons-gateway/utils/data/admins"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/teams"
	"github.com/IIP-Design/commons-gateway/utils/data/users"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
)

// updateAdminHandler handles the request to edit an existing admin user.
// It ensures that the required data is present before continuing on to
// update the team data.
func updateAdminHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	admin, err := data.ExtractAdminUser(event.Body)

	if err != nil {
		return msgs.SendServerError(err)
	}

	// Ensure that the user we intend to modify exists.
	_, adminExists, err := users.CheckForExistingAdminUser(admin.Email)

	if err != nil {
		logs.LogError(err, "Check For Admin Error")
		return msgs.SendServerError(err)
	} else if !adminExists {
		err = fmt.Errorf("the user %s has not been registered as an admin", admin.Email)

		logs.LogError(err, "Check For Admin Error")
		return msgs.SendCustomError(err, 404)
	}

	// Ensure that the user's assigned team exists.
	exists, err := teams.CheckForExistingTeamById(admin.Team)

	if err != nil {
		logs.LogError(err, "Check For Team Error")
		return msgs.SendServerError(err)
	} else if !exists {
		err = fmt.Errorf("no team with the id %s exists", admin.Team)

		logs.LogError(err, "Check For Team Error")
		return msgs.SendCustomError(err, 404)
	}

	err = admins.UpdateAdmin(admin)

	if err != nil {
		logs.LogError(err, "Update Admin Error")
		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(updateAdminHandler)
}
