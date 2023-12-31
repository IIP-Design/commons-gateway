package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/xid"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/users"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
	"github.com/IIP-Design/commons-gateway/utils/queue"
	"github.com/IIP-Design/commons-gateway/utils/randstr"
)

// registerMfaRequest saves the generated 2FA and request id to the database.
// These values are referenced when authenticating the user.
func registerMfaRequest(requestId xid.ID, code string) error {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	currentTime := time.Now()

	insertMfa := `INSERT INTO mfa( request_id, code, date_created ) VALUES ( $1, $2, $3 );`
	_, err = pool.Exec(insertMfa, requestId.String(), code, currentTime)

	if err != nil {
		logs.LogError(err, "Save MFA Request Query Error")
	}

	return err
}

// initiateEmailQueue sends the 2FA code to the SQS queue
// that manages the the sending of 2FA emails.
func initiateEmailQueue(username string, code string) error {
	var err error

	// Retrieve the user data.
	pool := data.ConnectToDB()
	defer pool.Close()

	var user data.User

	query := `SELECT email, first_name, last_name FROM guests WHERE email = $1;`
	err = pool.QueryRow(query, username).Scan(&user.Email, &user.NameFirst, &user.NameLast)

	if err != nil {
		logs.LogError(err, "Retrieve User Data Error")
		return err
	}

	// Prepare the message sent by SQS.
	body := map[string]any{
		"user": user,
		"code": code,
	}

	json, err := json.Marshal(body)

	if err != nil {
		logs.LogError(err, "Failed to Marshal SQS Body")
		return err
	}

	queueUrl := os.Getenv("EMAIL_MFA_QUEUE")

	// Send the message to SQS.
	messageId, err := queue.SendToQueue(string(json), queueUrl)

	if err != nil {
		logs.LogError(err, "Failed to Send Queue Message")
		return err
	}

	fmt.Println("Sent message with ID: " + messageId)

	return err
}

// generateMfaHandler creates a one-time code to be used as a second
// factor when authenticating guest users.
func generateMfaHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	username := event.QueryStringParameters["username"]

	if username == "" {
		logs.LogError(nil, "Missing Parameter Error - username")
		return msgs.SendCustomError(errors.New("user email not provided"), 400)
	}

	// Ensure that the user requesting a 2FA code exists.
	_, exists, err := users.CheckForExistingGuestUser(username)

	if err != nil {
		logs.LogError(err, "Check For User Error")
		return msgs.SendCustomError(errors.New("load guest error"), 500)
	} else if !exists {
		logs.LogError(fmt.Errorf("user %s not found", username), "Guest User Not Found Error")
		return msgs.SendCustomError(errors.New("no such user"), 404)
	}

	// Generate the 2FA code.
	requestId := xid.New()
	code, err := randstr.RandDigitBytes(6)

	if err != nil {
		logs.LogError(err, "Failed to Generate 2FA Code")
		return msgs.SendServerError(err)
	}

	// Save the 2FA request.
	err = registerMfaRequest(requestId, code)

	if err != nil {
		logs.LogError(err, "Failed to Register 2FA Code")
		return msgs.SendServerError(err)
	}

	// Email the user their code.
	err = initiateEmailQueue(username, code)

	if err != nil {
		logs.LogError(err, "Failed to Send 2FA Code")
		return msgs.SendCustomError(errors.New("internal error"), 500)
	}

	// Return the 2FA request id to the application.
	resp := map[string]any{
		"requestId": requestId,
	}

	body, err := msgs.MarshalBody(resp)

	if err != nil {
		logs.LogError(err, "Failed to Marshal Response Body")
		return msgs.SendServerError(err)
	}

	return msgs.PrepareResponse(body)
}

func main() {
	lambda.Start(generateMfaHandler)
}
