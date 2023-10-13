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
func handleTeamCreation(teamName string, aprimoName string) error {
	var err error

	exists, err := teams.CheckForExistingTeam(teamName)

	if err != nil {
		return err
	} else if exists {
		return errors.New("a team with this name already exists")
	}

	err = teams.CreateTeam(teamName, aprimoName)

	return err
}

// NewTeamHandler handles the request to add a new team for uploading. It
// ensures that the required data is present before continuing on to recording
// the team name and setting it to active.
func NewTeamHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	parsed, err := data.ParseBodyData(event.Body)

	team := parsed.TeamName
	aprimo_name := parsed.AprimoName

	if err != nil {
		return msgs.SendServerError(err)
	} else if team == "" {
		logs.LogError(nil, "Team name not provided in request.")
		return msgs.Response{StatusCode: 400}, errors.New("data missing from request")
	}

	err = handleTeamCreation(team, aprimo_name)

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
	lambda.Start(NewTeamHandler)
}
