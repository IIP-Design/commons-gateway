package main

import (
	"context"
	"errors"
	"time"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/users"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// deactivateGuest opens a database connection and sets the given guest's
// access expiration date to the current time.
func deactivateGuest(email string) error {
	pool := data.ConnectToDB()
	defer pool.Close()

	currentTime := time.Now()
	deactivatedTime := currentTime.Add(time.Duration(-1) * time.Minute)

	query := `UPDATE invites SET expiration = $1 WHERE invitee = $2`
	_, err := pool.Exec(query, deactivatedTime, email)

	if err != nil {
		logs.LogError(err, "Deactivate Guest Query Error")
	}

	return err
}

// guestDeactivateHandler handles the request to edit an existing guest user.
// It ensures that the required data is present before continuing on to
// update the team data.
func guestDeactivateHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	id := event.QueryStringParameters["id"]

	if id == "" {
		return msgs.SendServerError(errors.New("user id not provided"))
	}

	// Ensure that the user we intend to modify exists.
	_, exists, err := users.CheckForExistingUser(id, "guests")

	if !exists {
		return msgs.SendCustomError(errors.New("user does not exist"), 404)
	} else if err != nil {
		return msgs.SendServerError(err)
	}

	err = deactivateGuest(id)

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(guestDeactivateHandler)
}
