package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/IIP-Design/commons-gateway/utils/data/guests"
	"github.com/IIP-Design/commons-gateway/utils/data/users"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
)

// getGuestHandler handles the request to retrieve a single admin user based on email address.
func getGuestHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	id := event.QueryStringParameters["id"]

	if id == "" {
		return msgs.SendCustomError(errors.New("user id not provided"), 400)
	}

	// Ensure the user exists doesn't already have access.
	_, exists, err := users.CheckForExistingGuestUser(id)

	if err != nil {
		logs.LogError(err, "Check For Guest User Error")
		return msgs.SendServerError(err)
	} else if !exists {
		err = fmt.Errorf("%s is not registered as a guest user", id)

		logs.LogError(err, "Guest User Not Found Error")
		return msgs.SendCustomError(errors.New("user does not exist"), 404)
	}

	guest, err := guests.RetrieveGuest(id)

	if err != nil {
		logs.LogError(err, "Retrieve Guest Error")
		return msgs.SendServerError(err)
	}

	body, err := msgs.MarshalBody(guest)

	if err != nil {
		logs.LogError(err, "Marshal Body Error")
		return msgs.SendServerError(err)
	}

	return msgs.PrepareResponse(body)
}

func main() {
	lambda.Start(getGuestHandler)
}
