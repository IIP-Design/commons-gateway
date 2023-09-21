package main

import (
	"context"
	"errors"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/guests"
	"github.com/IIP-Design/commons-gateway/utils/data/teams"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
	"github.com/IIP-Design/commons-gateway/utils/security/jwt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// GuestUpdateHandler handles the request to edit an existing guest user.
// It ensures that the required data is present before continuing on to
// update the team data.
func GuestUpdateHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	code, err := jwt.RequestIsAuthorized(event, []string{"super admin", "admin", "guest admin"})
	if err != nil {
		return msgs.SendAuthError(err, code)
	}

	guest, err := data.ExtractGuestUser(event.Body)

	if err != nil {
		return msgs.SendServerError(err)
	}

	// Ensure that the user we intend to modify exists.
	_, userExists, err := data.CheckForExistingUser(guest.Email, "guests")

	if err != nil {
		return msgs.SendServerError(err)
	} else if !userExists {
		return msgs.SendServerError(errors.New("this user has not been registered"))
	}

	// Ensure that the user's assigned team exists.
	exists, err := teams.CheckForExistingTeamById(guest.Team)

	if err != nil {
		return msgs.SendServerError(err)
	} else if !exists {
		return msgs.SendServerError(errors.New("no team with the provided id exists"))
	}

	err = guests.UpdateGuest(guest)

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(GuestUpdateHandler)
}
