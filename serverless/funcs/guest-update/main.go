package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/guests"
	"github.com/IIP-Design/commons-gateway/utils/data/teams"
	"github.com/IIP-Design/commons-gateway/utils/data/users"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
)

// GuestUpdateHandler handles the request to edit an existing guest user.
// It ensures that the required data is present before continuing on to
// update the team data.
func GuestUpdateHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	guest, err := data.ExtractGuestUser(event.Body)

	if err != nil {
		return msgs.SendServerError(err)
	}

	// Ensure that the user we intend to modify exists.
	_, userExists, err := users.CheckForExistingUser(guest.Email, "guests")

	if err != nil {
		logs.LogError(err, "Check For User Error")
		return msgs.SendServerError(err)
	} else if !userExists {
		logs.LogError(fmt.Errorf("user %s not found", guest.Email), "User Not Found Error")
		return msgs.SendServerError(errors.New("this user has not been registered"))
	}

	// Ensure that the user's assigned team exists.
	exists, err := teams.CheckForExistingTeamById(guest.Team)

	if err != nil {
		logs.LogError(err, "Check For Team Error")
		return msgs.SendServerError(err)
	} else if !exists {
		logs.LogError(fmt.Errorf("team with id %s not found", guest.Team), "Team Not Found Error")
		return msgs.SendServerError(errors.New("no team with the provided id exists"))
	}

	err = guests.UpdateGuest(guest)

	if err != nil {
		logs.LogError(err, "Update Guest Error")
		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(GuestUpdateHandler)
}
