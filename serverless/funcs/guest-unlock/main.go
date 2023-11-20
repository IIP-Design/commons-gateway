package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IIP-Design/commons-gateway/utils/data/creds"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// unlockGuestAccount resets a given guest's account back to unlocked.
func unlockGuestAccount(guest string) error {
	pool := data.ConnectToDB()
	defer pool.Close()

	query := `UPDATE guests SET locked = false WHERE email = $1;`
	_, err := pool.Exec(query, guest)

	if err != nil {
		logs.LogError(err, "Unlock Account Query Error")
	}

	return err
}

// unlockGuestHandler manages SQS messages to set a locked guest user's account back to unlocked status.
func unlockGuestHandler(ctx context.Context, event events.SQSEvent) (msgs.Response, error) {
	for _, message := range event.Records {
		eventMessageId := message.MessageId
		body := message.Body

		var eventData data.GuestUnlockInitEvent

		err := json.Unmarshal([]byte(body), &eventData)

		if err != nil {
			logs.LogError(err, fmt.Sprintf("Unable to unmarshal body of message %s", eventMessageId))
			return msgs.SendServerError(err)
		}

		err = unlockGuestAccount(eventData.Username)

		if err != nil {
			return msgs.SendServerError(err)
		}

		// Reset the login counter
		err = creds.ClearUnsuccessfulLoginAttempts(eventData.Username)

		if err != nil {
			return msgs.SendServerError(err)
		}
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(unlockGuestHandler)
}
