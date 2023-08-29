package main

import (
	"context"
	"errors"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/teams"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// TeamUpdateHandler handles the request to edit an existing team. It
// ensures that the required data is present before continuing on to
// update the team data.
func TeamUpdateHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	parsed, err := data.ParseBodyData(event.Body)

	active := parsed.Active
	name := parsed.TeamName
	team := parsed.TeamId

	if err != nil {
		return msgs.SendServerError(err)
	} else if team == "" {
		logs.LogError(nil, "Team data not provided in request.")
		return msgs.Response{StatusCode: 400}, errors.New("data missing from request")
	}

	exists, err := teams.CheckForExistingTeamById(team)

	if err != nil {
		return msgs.SendServerError(err)
	} else if !exists {
		return msgs.SendServerError(errors.New("no team with this id exists"))
	}

	if name != "" {
		// If both active status and team name provided update full team info.
		err = teams.UpdateTeam(team, name, active)
	} else {
		// If only status provided, update status.
		err = teams.UpdateTeamStatus(team, active)
	}

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(TeamUpdateHandler)
}
