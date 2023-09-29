package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IIP-Design/commons-gateway/utils/aprimo"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RecordCreationResponse struct {
	Id string `json:"id"`
}

func ParseEventBody(body string) (aprimo.FileRecordInitEvent, error) {
	var parsed aprimo.FileRecordInitEvent

	b := []byte(body)
	err := json.Unmarshal(b, &parsed)

	return parsed, err
}

func submitRecord(description string, event aprimo.FileRecordInitEvent, team string, token string) (string, error) {
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
		},
		"files": {
			"master": "%s",
			"addOrUpdate": [
				{
					"versions": {
						"addOrUpdate": [
							{
								"id": "%s",
								"filename": "%s"
							}
						]
					}
				}
			]
		}
	}`, description, event.Key, team, event.FileToken, event.FileToken, event.Key)

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

// CreateAprimoRecord initiates the creation of a new record in Aprimo.
func CreateAprimoRecord(ctx context.Context, event events.SQSEvent) {
	token, err := aprimo.GetAuthToken()

	if err != nil {
		logs.LogError(err, "Unable to Authenticate Error")

		return
	}

	pool := data.ConnectToDB()
	defer pool.Close()

	for _, record := range event.Records {
		fileInfo, err := ParseEventBody(record.Body)
		if err != nil {
			logs.LogError(err, "SQS event body parse error")
			return
		}

		var description string
		var team string

		query := "SELECT uploads.description, teams.aprimo_name FROM uploads INNER JOIN teams ON uploads.team_id=teams.id WHERE uploads.s3_id = $1"
		err = pool.QueryRow(query, fileInfo.Key).Scan(&description, &team)

		if err != nil {
			logs.LogError(err, "Retrieve Upload Metadata Query Error")
			return
		}

		recordId, err := submitRecord(description, fileInfo, team, token)
		if err != nil {
			logs.LogError(err, "Aprimo Record Create Error")
		} else {
			log.Println(recordId)
		}

		query = "UPDATE uploads SET aprimo_record_id = $1 WHERE s3_id = $2"
		_, err = pool.Exec(query, recordId, fileInfo.Key)
		if err != nil {
			logs.LogError(err, "Aprimo Record ID Save Error")
		}
	}
}

func main() {
	lambda.Start(CreateAprimoRecord)
}
