package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/IIP-Design/commons-gateway/utils/aprimo"
	"github.com/IIP-Design/commons-gateway/utils/logs"
)

func ParseEventBody(body string) (aprimo.FileRecordUpdateEvent, error) {
	var parsed aprimo.FileRecordUpdateEvent

	b := []byte(body)
	err := json.Unmarshal(b, &parsed)

	return parsed, err
}

func updateAprimoRecord(ctx context.Context, event events.SQSEvent) error {
	var err error

	// Retrieve Aprimo auth token
	token, err := aprimo.GetAuthToken()

	if err != nil {
		logs.LogError(err, "Unable to Authenticate Error")
		return err
	}

	for _, message := range event.Records {
		fileInfo, err := ParseEventBody(message.Body)
		if err != nil {
			logs.LogError(err, "Failed to Unmarshal Body")
			return err
		}

	}

	return err
}

func main() {
	lambda.Start(updateAprimoRecord)
}
