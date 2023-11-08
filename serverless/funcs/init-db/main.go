package main

import (
	"context"

	initdb "github.com/IIP-Design/commons-gateway/utils/init"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/lambda"
)

// initDBHandler handles the request to set up the database.
func initDBHandler(ctx context.Context) (msgs.Response, error) {
	migExists := initdb.CheckForTable("migrations")

	// If the migrations table does not exist run the initial setup
	if !migExists {
		err := initdb.InitializeDatabase()

		if err != nil {
			return msgs.SendServerError(err)
		}
	}

	err := initdb.ApplyMigrations()

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(initDBHandler)
}
