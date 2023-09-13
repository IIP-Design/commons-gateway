package main

import (
	"context"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/guests"
	"github.com/IIP-Design/commons-gateway/utils/jwt"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// GetGuestsHandler handles the request to retrieve a list of guest users.
// If a 'team' argument is provided in the body of the request it will filter
// the response to show only the guests assigned to that team.
func GetGuestsHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	code, err := jwt.RequestIsAuthorized(event, []string{"super admin", "admin"})
	if err != nil {
		return msgs.SendAuthError(err, code)
	}

	parsed, err := data.ParseBodyData(event.Body)

	team := parsed.TeamId
	role := parsed.Role

	if err != nil {
		return msgs.SendServerError(err)
	}

	guests, err := guests.RetrieveGuests(team, role)

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
