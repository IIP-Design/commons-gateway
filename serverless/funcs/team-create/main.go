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

// handleTeamCreation coordinates all the actions associated with creating a new team.
func handleTeamCreation(teamName string, aprimoName string) (bool, error) {
	var err error
	var exists bool

	exists, err = teams.CheckForExistingTeam(teamName)

	if err != nil || exists {
		return exists, err
	}

	err = teams.CreateTeam(teamName, aprimoName)

	return exists, err
}

// newTeamHandler handles the request to add a new team for uploading. It
// ensures that the required data is present before continuing on to recording
// the team name and setting it to active.
func newTeamHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	parsed, err := data.ParseBodyData(event.Body)

	team := parsed.TeamName
	aprimo_name := parsed.AprimoName

	if err != nil {
		return msgs.SendServerError(err)
	} else if team == "" {
		err := errors.New("data missing from request")
		logs.LogError(err, "Team name not provided in request.")
		return msgs.SendCustomError(err, 400)
	}

	exists, err := handleTeamCreation(team, aprimo_name)

	if exists {
		return msgs.SendCustomError(errors.New("a team with this name already exists"), 409)
	} else if err != nil {
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
	lambda.Start(newTeamHandler)
}
