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
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/rs/xid"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
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
	_, err = pool.Exec(insertMfa, requestId, code, currentTime)

	if err != nil {
		logs.LogError(err, "Save MFA Request Query Error")
	}

	return err
}

// initiateEmailQueue sends the 2FA code to the SQS queue
// that manages the the sending of 2FA emails.
func initiateEmailQueue(username string, code string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		logs.LogError(err, "Error Loading AWS Config")
		return err
	}

	client := sqs.NewFromConfig(cfg)

	queue := os.Getenv("EMAIL_QUEUE")

	body := map[string]string{
		"email": username,
		"code":  code,
	}

	json, err := json.Marshal(body)

	if err != nil {
		logs.LogError(err, "Failed to Marshal SQS Body")
		return err
	}

	messageInput := &sqs.SendMessageInput{
		DelaySeconds: 0,
		MessageBody:  aws.String(string(json)),
		QueueUrl:     &queue,
	}

	resp, err := client.SendMessage(context.TODO(), messageInput)

	if err != nil {
		logs.LogError(err, "Failed to Send Queue Message")
		return err
	}

	fmt.Println("Sent message with ID: " + *resp.MessageId)

	return err
}

// generateMfaHandler creates a one-time code to be used as a second
// factor when authenticating guest users.
func generateMfaHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	username := event.QueryStringParameters["username"]

	if username == "" {
		return msgs.SendServerError(errors.New("user email not provided"))
	}

	// Ensure that the user requesting a 2FA code exists.
	_, exists, err := data.CheckForExistingUser(username, "guests")

	if err != nil || !exists {
		return msgs.SendServerError(errors.New("user does not exist"))
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
		logs.LogError(err, "Failed to Generate 2FA Code")
		return msgs.SendServerError(err)
	}

	// Email the user their code.
	err = initiateEmailQueue(username, code)

	if err != nil {
		logs.LogError(err, "Failed to Send 2FA Code")
		return msgs.SendAuthError(errors.New("internal error"), 500)
	}

	// Return the 2FA request id to the application.
	resp := map[string]any{
		"requestId": requestId,
	}

	body, err := msgs.MarshalBody(resp)

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.PrepareResponse(body)
}

func main() {
	lambda.Start(generateMfaHandler)
}
