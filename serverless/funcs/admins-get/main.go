package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/IIP-Design/commons-gateway/utils/data/admins"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
)

// GetAdminsHandler handles the request to retrieve a list of all admin users.
func GetAdminsHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	admins, err := admins.RetrieveAdmins()

	if err != nil {
		return msgs.SendServerError(err)
	}

	body, err := msgs.MarshalBody(admins)

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.PrepareResponse(body)
}

func main() {
	lambda.Start(GetAdminsHandler)
}
