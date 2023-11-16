package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/guests"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
)

// GetUploaderHandler handles the request to retrieve a list of uploader guests on the guest admin's team.
func GetUploaderHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	parsed, err := data.ParseBodyData(event.Body)

	team := parsed.TeamId

	if err != nil {
		logs.LogError(err, "Parsing Body Error")
		return msgs.SendServerError(err)
	}

	guests, err := guests.RetrieveUploaders(team)

	if err != nil {
		logs.LogError(err, "Uploaders Retrieve Error")
		return msgs.SendServerError(err)
	}

	body, err := msgs.MarshalBody(guests)

	if err != nil {
		logs.LogError(err, "Failed to Unmarshal Body")
		return msgs.SendServerError(err)
	}

	return msgs.PrepareResponse(body)
}

func main() {
	lambda.Start(GetUploaderHandler)
}
