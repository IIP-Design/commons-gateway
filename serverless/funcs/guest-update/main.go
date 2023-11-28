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

// guestUpdateHandler handles the request to edit an existing guest user.
// It ensures that the required data is present before continuing on to
// update the team data.
func guestUpdateHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	guest, err := data.ExtractGuestUser(event.Body)

	if err != nil {
		return msgs.SendServerError(err)
	}

	// Ensure that the user we intend to modify exists.
	_, userExists, err := users.CheckForExistingGuestUser(guest.Email)

	if err != nil {
		logs.LogError(err, "Check For Guest User Error")
		return msgs.SendServerError(err)
	} else if !userExists {
		err = fmt.Errorf("user %s is not registered as a guest", guest.Email)

		logs.LogError(err, "User Not Found Error")
		return msgs.SendCustomError(errors.New("this user has not been registered"), 404)
	}

	// Ensure that the user's assigned team exists.
	exists, err := teams.CheckForExistingTeamById(guest.Team)

	if err != nil {
		logs.LogError(err, "Check For Team Error")
		return msgs.SendServerError(err)
	} else if !exists {
		logs.LogError(fmt.Errorf("team with id %s not found", guest.Team), "Team Not Found Error")
		return msgs.SendCustomError(errors.New("no team with the provided id exists"), 404)
	}

	err = guests.UpdateGuest(guest)

	if err != nil {
		logs.LogError(err, "Update Guest Error")
		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(guestUpdateHandler)
}
