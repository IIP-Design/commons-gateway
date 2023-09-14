package main

import (
	"context"
	"errors"
	"time"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
	"github.com/IIP-Design/commons-gateway/utils/security/jwt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func deactivateGuest(email string) error {
	pool := data.ConnectToDB()
	defer pool.Close()

	currentTime := time.Now()
	deactivatedTime := currentTime.Add(time.Duration(-1) * time.Minute)

	query := `UPDATE guests SET expiration = $1, date_modified = $2 WHERE email = $3`
	_, err := pool.Exec(query, deactivatedTime, currentTime, email)

	if err != nil {
		logs.LogError(err, "Deactivate Guest Query Error")
	}

	return err
}

// GuestDeactivateHandler handles the request to edit an existing guest user.
// It ensures that the required data is present before continuing on to
// update the team data.
func GuestDeactivateHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	code, err := jwt.RequestIsAuthorized(event, []string{"super admin", "admin"})
	if err != nil {
		return msgs.SendAuthError(err, code)
	}

	id := event.QueryStringParameters["id"]

	if id == "" {
		return msgs.SendServerError(errors.New("user id not provided"))
	}

	// Ensure that the user we intend to modify exists.
	exists, err := data.CheckForExistingUser(id, "guests")

	if err != nil || !exists {
		return msgs.SendServerError(errors.New("user does not exist"))
	}

	err = deactivateGuest(id)

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(GuestDeactivateHandler)
}
