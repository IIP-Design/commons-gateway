package main

import (
	"context"

	"github.com/IIP-Design/commons-gateway/utils/data/teams"
	"github.com/IIP-Design/commons-gateway/utils/jwt"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// GetTeamsHandler handles the request to retrieve a list of all the teams.
func GetTeamsHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	code, err := jwt.RequestIsAuthorized(event, []string{"super admin", "admin", "guest"})
	if err != nil {
		return msgs.SendAuthError(err, code)
	}

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
	lambda.Start(GetTeamsHandler)
}
