package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/IIP-Design/commons-gateway/utils/aprimo"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/IIP-Design/commons-gateway/utils/queue"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RecordCreationResponse struct {
	Id string `json:"id"`
}

func submitRecord(description string, key string, team string, token string) (string, error) {
	var id string
	var err error

	reqBody := fmt.Sprintf(`{
		"status":"draft",
		"fields": {
			"addOrUpdate": [
				{
					"Name": "Description",
					"localizedValues": [
						{ "value": "%s" }
					]
				},
				{
					"Name": "DisplayTitle",
					"localizedValues": [
						{ "value": "%s" }
					]
				},
				{
					"Name": "Team",
					"localizedValues": [
						{ "values": ["%s"] }
					]
			}
			]
		}
	}`, description, key, team)

	respBody, _, err := aprimo.PostJsonData("records", token, reqBody)
	if err != nil {
		return id, err
	}

	var res RecordCreationResponse
	err = json.Unmarshal(respBody, &res)
	if err != nil {
		return id, err
	}

	return res.Id, nil
}

func SendRecordEvent(aprimoId string, filename string, fileType string) (string, error) {
	var messageId string
	var err error

	event := aprimo.FileRecordInitEvent{
		AprimoId: aprimoId,
		Filename: filename,
		FileType: fileType,
	}

	json, err := json.Marshal(event)

	if err != nil {
		logs.LogError(err, "Failed to Marshal SQS Body")
		return messageId, err
	}

	queueUrl := os.Getenv("RECORD_INIT_QUEUE")

	// Send the message to SQS.
	return queue.SendToQueue(string(json), queueUrl)
}

// CreateAprimoRecord initiates the creation of a new record in Aprimo.
func CreateAprimoRecord(ctx context.Context, event events.S3Event) {
	token, err := aprimo.GetAuthToken()

	if err != nil {
		logs.LogError(err, "Unable to Authenticate Error")

		return
	}

	pool := data.ConnectToDB()
	defer pool.Close()

	for _, record := range event.Records {
		key := record.S3.Object.Key
		var description string
		var fileType string
		var team string

		query := "SELECT uploads.description, uploads.file_type, teams.aprimo_name FROM uploads INNER JOIN teams ON uploads.team_id=teams.id WHERE uploads.s3_id = $1"
		err = pool.QueryRow(query, key).Scan(&description, &fileType, &team)

		if err != nil {
			logs.LogError(err, "Retrieve Upload Metadata Query Error")

			return
		}

		recordId, err := submitRecord(description, key, team, token)
		if err != nil {
			logs.LogError(err, "Aprimo Record Create Error")
		} else {
			log.Println(recordId)
		}

		query = "UPDATE uploads SET aprimo_record_id = $1 WHERE s3_id = $2"
		_, err = pool.Exec(query, recordId, key)
		if err != nil {
			logs.LogError(err, "Aprimo Record ID Save Error")
		}

		// Pass along to file upload
		messageId, err := SendRecordEvent(recordId, key, fileType)
		if err != nil {
			logs.LogError(err, "send record update event error")
		} else {
			log.Println(messageId)
		}
	}
}

func main() {
	lambda.Start(CreateAprimoRecord)
}
