package main

import (
	"context"

	data "github.com/IIP-Design/commons-gateway/utils/data"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/lambda"
)

// GetTeamsHandler handles the request to retrieve a list of all the teams.
func GetTeamsHandler(ctx context.Context) (msgs.Response, error) {
	var err error

	teams, err := data.RetrieveTeams()

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
	lambda.Start(GetTeamsHandler)
}
