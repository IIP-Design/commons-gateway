package main

import (
	"context"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/guests"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// GetGuestsHandler handles the request to retrieve a list of guest users.
// If a 'team' argument is provided in the body of the request it will filter
// the response to show only the guests assigned to that team.
func GetGuestsHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	var err error

	parsed, err := data.ParseBodyData(event.Body)

	team := parsed.TeamId

	if err != nil {
		return msgs.SendServerError(err)
	}

	guests, err := guests.RetrieveGuests(team)

	if err != nil {
		return msgs.SendServerError(err)
	}

	body, err := msgs.MarshalBody(guests)

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.PrepareResponse(body)
}

func main() {
	lambda.Start(GetGuestsHandler)
}
