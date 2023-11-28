package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/IIP-Design/commons-gateway/utils/data/admins"
	"github.com/IIP-Design/commons-gateway/utils/data/users"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
)

// getAdminHandler handles the request to retrieve a single admin user based on email address.
func getAdminHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	username := event.QueryStringParameters["username"]

	if username == "" {
		return msgs.SendServerError(errors.New("user name not provided"))
	}

	// Ensure the user exists and already has access.
	_, exists, err := users.CheckForExistingAdminUser(username)

	if err != nil {
		logs.LogError(err, "Check For Admin Error")
		return msgs.SendServerError(err)
	} else if !exists {
		err := fmt.Errorf("user %s does not exist", username)

		logs.LogError(err, "Check For Admin Error")
		return msgs.SendCustomError(err, 400)
	}

	admin, err := admins.RetrieveAdmin(username)

	if err != nil {
		logs.LogError(err, "Retrieve Admin Error")
		return msgs.SendServerError(err)
	}

	body, err := msgs.MarshalBody(admin)

	if err != nil {
		logs.LogError(err, "Marshal Body Error")
		return msgs.SendServerError(err)
	}

	return msgs.PrepareResponse(body)
}

func main() {
	lambda.Start(getAdminHandler)
}
