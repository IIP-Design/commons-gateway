package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/IIP-Design/commons-gateway/utils/data/admins"
	"github.com/IIP-Design/commons-gateway/utils/data/users"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
)

// getAdminHandler handles the request to retrieve a single admin user based on email address.
func getAdminHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	username := event.QueryStringParameters["username"]

	if username == "" {
		return msgs.SendServerError(errors.New("user name not provided"))
	}

	// Ensure the user exists and already has access.
	_, exists, err := users.CheckForExistingUser(username, "admins")

	if !exists {
		return msgs.SendCustomError(errors.New("user is not an admin"), 404)
	} else if err != nil {
		return msgs.SendServerError(err)
	}

	admin, err := admins.RetrieveAdmin(username)

	if err != nil {
		return msgs.SendServerError(err)
	}

	body, err := msgs.MarshalBody(admin)

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.PrepareResponse(body)
}

func main() {
	lambda.Start(getAdminHandler)
}
