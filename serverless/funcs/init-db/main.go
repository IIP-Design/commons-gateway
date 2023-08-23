package main

import (
	"context"

	initdb "github.com/IIP-Design/commons-gateway/utils/data/init"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/lambda"
)

// InitDBHandler handles the request to retrieve a list of all the teams.
func InitDBHandler(ctx context.Context) (msgs.Response, error) {
	err := initdb.InitializeDatabase()

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.PrepareResponse([]byte("success"))
}

func main() {
	lambda.Start(InitDBHandler)
}
