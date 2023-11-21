package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/IIP-Design/commons-gateway/utils/data/teams"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
)

// getTeamsHandler handles the request to retrieve a list of all the teams.
func getTeamsHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
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
	lambda.Start(getTeamsHandler)
}
