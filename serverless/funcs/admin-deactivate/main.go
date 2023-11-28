package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/users"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
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

// deactivateAdminHandler handles the request to deactivate an existing admin.
// It ensures that the required data is present before continuing on to update the admin data.
func deactivateAdminHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	username := event.QueryStringParameters["username"]

	if username == "" {
		return msgs.SendServerError(errors.New("admin id not provided"))
	}

	// Ensure that the user we intend to modify exists.
	_, exists, err := users.CheckForExistingAdminUser(username)

	if !exists {
		err = fmt.Errorf("%s does not exist as an admin user", username)

		logs.LogError(err, "Deactivate Admin Error")
		return msgs.SendCustomError(err, 404)
	} else if err != nil {
		logs.LogError(err, "Deactivate Admin Error")
		return msgs.SendServerError(err)
	}

	err = deactivateAdmin(username)

	if err != nil {
		logs.LogError(err, "Deactivate Admin Error")
		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(deactivateAdminHandler)
}
