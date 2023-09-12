package main

import (
	"context"
	"errors"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/guests"
	"github.com/IIP-Design/commons-gateway/utils/jwt"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// GetGuestHandler handles the request to retrieve a single admin user based on email address.
func GetGuestHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	_, err := jwt.RequestIsAuthorized(event, []string{"super admin", "admin", "guest"})
	if err != nil {
		return msgs.SendServerError(err)
	}

	id := event.QueryStringParameters["id"]

	if id == "" {
		return msgs.SendServerError(errors.New("user id not provided"))
	}

	// Ensure the user exists doesn't already have access.
	exists, err := data.CheckForExistingUser(id, "guests")

	if err != nil || !exists {
		return msgs.SendServerError(errors.New("user does not exist"))
	}

	guest, err := guests.RetrieveGuest(id)

	if err != nil {
		return msgs.SendServerError(err)
	}

	body, err := msgs.MarshalBody(guest)

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.PrepareResponse(body)
}

func main() {
	lambda.Start(GetGuestHandler)
}
