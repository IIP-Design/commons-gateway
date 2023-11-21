package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/IIP-Design/commons-gateway/utils/data/guests"
	"github.com/IIP-Design/commons-gateway/utils/data/users"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
)

// getGuestHandler handles the request to retrieve a single admin user based on email address.
func getGuestHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	id := event.QueryStringParameters["id"]

	if id == "" {
		return msgs.SendCustomError(errors.New("user id not provided"), 400)
	}

	// Ensure the user exists doesn't already have access.
	_, exists, err := users.CheckForExistingUser(id, "guests")

	if !exists {
		return msgs.SendCustomError(errors.New("user does not exist"), 404)
	} else if err != nil {
		return msgs.SendServerError(err)
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
	lambda.Start(getGuestHandler)
}
