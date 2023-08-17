package main

import (
	"context"

	data "github.com/IIP-Design/commons-gateway/utils/data"
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

	team := parsed.Team

	if err != nil {
		return msgs.Response{StatusCode: 500}, err
	}

	guests, err := data.RetrieveGuests(team)

	if err != nil {
		return msgs.Response{StatusCode: 500}, err
	}

	body, err := msgs.MarshalBody(guests)

	if err != nil {
		return msgs.Response{StatusCode: 500}, err
	}

	return msgs.PrepareResponse(body)
}

func main() {
	lambda.Start(GetGuestsHandler)
}
