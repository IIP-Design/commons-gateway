package main

import (
	"context"

	data "github.com/IIP-Design/commons-gateway/utils/data"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/lambda"
)

// GetAdminsHandler handles the request to create a new administrative user. It
// ensures that the required data is present before continuing on to recording
// the user's email in the list of admins.
func GetAdminsHandler(ctx context.Context) (msgs.Response, error) {
	var err error

	admins, err := data.RetrieveAdmins()

	if err != nil {
		return msgs.Response{StatusCode: 500}, err
	}

	body, err := msgs.MarshalBody(admins)

	if err != nil {
		return msgs.Response{StatusCode: 500}, err
	}

	return msgs.PrepareResponse(body)
}

func main() {
	lambda.Start(GetAdminsHandler)
}
