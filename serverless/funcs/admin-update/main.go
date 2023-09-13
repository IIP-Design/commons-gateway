package main

import (
	"context"
	"errors"

	"github.com/IIP-Design/commons-gateway/utils/data/admins"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/teams"
	"github.com/IIP-Design/commons-gateway/utils/jwt"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// UpdateAdminHandler handles the request to edit an existing admin user.
// It ensures that the required data is present before continuing on to
// update the team data.
func UpdateAdminHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	code, err := jwt.RequestIsAuthorized(event, []string{"super admin"})
	if err != nil {
		return msgs.SendAuthError(err, code)
	}

	admin, err := data.ExtractAdminUser(event.Body)

	if err != nil {
		return msgs.SendServerError(err)
	}

	// Ensure that the user we intend to modify exists.
	adminExists, err := data.CheckForExistingUser(admin.Email, "admins")

	if err != nil {
		return msgs.SendServerError(err)
	} else if !adminExists {
		return msgs.SendServerError(errors.New("this admin has not been registered"))
	}

	// Ensure that the user's assigned team exists.
	exists, err := teams.CheckForExistingTeamById(admin.Team)

	if err != nil {
		return msgs.SendServerError(err)
	} else if !exists {
		return msgs.SendServerError(errors.New("no team with the provided id exists"))
	}

	err = admins.UpdateAdmin(admin)

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(UpdateAdminHandler)
}
