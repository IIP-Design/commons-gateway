package main

import (
	"context"
	"errors"
	"time"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// deactivateAdmin sets an existing admin's `active` status to `false`.
func deactivateAdmin(email string) error {
	pool := data.ConnectToDB()
	defer pool.Close()

	currentTime := time.Now()

	query := `UPDATE admins SET active = false, date_modified = $1 WHERE email = $2`
	_, err := pool.Exec(query, currentTime, email)

	if err != nil {
		logs.LogError(err, "Deactivate Admin Query Error")
	}

	return err
}

// DeactivateAdminHandler handles the request to deactivate an existing admin.
// It ensures that the required data is present before continuing on to update the admin data.
func DeactivateAdminHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	username := event.QueryStringParameters["username"]

	if username == "" {
		return msgs.SendServerError(errors.New("admin id not provided"))
	}

	// Ensure that the user we intend to modify exists.
	exists, err := data.CheckForExistingUser(username, "admins")

	if err != nil || !exists {
		return msgs.SendServerError(errors.New("admin does not exist"))
	}

	err = deactivateAdmin(username)

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(DeactivateAdminHandler)
}
