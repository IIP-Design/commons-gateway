package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/teams"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
)

// TeamUpdateHandler handles the request to edit an existing team. It
// ensures that the required data is present before continuing on to
// update the team data.
func TeamUpdateHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	parsed, err := data.ParseBodyData(event.Body)

	active := parsed.Active
	name := parsed.TeamName
	team := parsed.TeamId
	aprimo_name := parsed.AprimoName

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
		err = teams.UpdateTeam(team, name, aprimo_name, active)
	} else {
		// If only status provided, update status.
		err = teams.UpdateTeamStatus(team, active)
	}

	if err != nil {
		return msgs.SendServerError(err)
	}

	// Return the full list of teams in the response.
	teams, err := teams.RetrieveTeams()

	if err != nil {
		return msgs.SendServerError(err)
	}

	body, err := msgs.MarshalBody(teams)

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.PrepareResponse(body)
}

func main() {
	lambda.Start(TeamUpdateHandler)
}
