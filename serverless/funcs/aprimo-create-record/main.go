package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	"github.com/IIP-Design/commons-gateway/utils/aprimo"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func ParseEventBody(body string) (aprimo.FileRecordInitEvent, error) {
	var parsed aprimo.FileRecordInitEvent

	b := []byte(body)
	err := json.Unmarshal(b, &parsed)

	return parsed, err
}

// handleCreateAprimoRecord initiates the creation of a new record in Aprimo.
func handleCreateAprimoRecord(ctx context.Context, event events.SQSEvent) error {
	token, err := aprimo.GetAuthToken()

	if err != nil {
		logs.LogError(err, "Unable to Authenticate Error")
		return err
	}

	log.Printf("SQS Events: %d\n", len(event.Records))

	pool := data.ConnectToDB()
	defer pool.Close()

	for _, record := range event.Records {
		fileInfo, err := ParseEventBody(record.Body)
		if err != nil {
			logs.LogError(err, "SQS event body parse error")
			return err
		}
		log.Printf("Event body: %s\n", record.Body)

		var description string
		var team string
		var aprimoRecordId sql.NullString

		query := "SELECT uploads.description, teams.aprimo_name, uploads.aprimo_record_id FROM uploads INNER JOIN teams ON uploads.team_id=teams.id WHERE uploads.s3_id = $1"
		err = pool.QueryRow(query, fileInfo.Key).Scan(&description, &team, &aprimoRecordId)

		if err != nil {
			logs.LogError(err, "Retrieve Upload Metadata Query Error")
			return err
		}

		if !aprimoRecordId.Valid { // No Aprimo record ID means this is likely a new record that we need to create
			recordId, err := aprimo.SubmitRecord(description, fileInfo, team, token)
			if err != nil {
				logs.LogError(err, "Aprimo Record Create Error")
				return err
			} else {
				log.Println(recordId)
			}

			query = "UPDATE uploads SET aprimo_record_id = $1, aprimo_record_dt = NOW() WHERE s3_id = $2"
			_, err = pool.Exec(query, recordId, fileInfo.Key)
			if err != nil {
				logs.LogError(err, "Aprimo Record ID Save Error")
				return err
			}
		} else { // If there's already an Aprimo ID, this is a replayed event and we don't want to act on it
			log.Printf("Object %s already has a record (%s), but the event was not deleted", fileInfo.Key, aprimoRecordId.String)
		}
	}

	return nil
}

func main() {
	lambda.Start(handleCreateAprimoRecord)
}
