package main

import (
	"context"
	"errors"

	"github.com/IIP-Design/commons-gateway/utils/data/admins"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// GetAdminHandler handles the request to retrieve a single admin user based on email address.
func GetAdminHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	parsed, err := data.ParseBodyData(event.Body)

	if err != nil {
		return msgs.SendServerError(err)
	}

	username := parsed.Username

	active, err := admins.CheckForActiveAdmin(username)

	if err != nil || !active {
		return msgs.SendServerError(errors.New("user not and active admin"))
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
	lambda.Start(GetAdminHandler)
}
