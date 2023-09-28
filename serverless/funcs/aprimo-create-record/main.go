package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/IIP-Design/commons-gateway/utils/aprimo"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// CreateAprimoRecord initiates the creation of a new record in Aprimo.
func CreateAprimoRecord(ctx context.Context, event events.S3Event) {
	token, err := aprimo.GetAuthToken()

	if err != nil {
		logs.LogError(err, "Unable to Authenticate Error")

		return
	}

	pool := data.ConnectToDB()
	defer pool.Close()

	client := &http.Client{}

	for _, record := range event.Records {
		key := record.S3.Object.Key
		var description string
		var team string

		query := "SELECT uploads.description, teams.aprimo_name FROM uploads INNER JOIN teams ON uploads.team_id=teams.id WHERE uploads.s3_id = $1"
		err = pool.QueryRow(query, key).Scan(&description, &team)

		if err != nil {
			logs.LogError(err, "Retrieve Upload Metadata Query Error")

			return
		}

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

		endpoint := aprimo.GetEndpointURL("records", false)

		jsonData := []byte(reqBody)
		bodyReader := bytes.NewReader(jsonData)

		request, err := http.NewRequest(
			http.MethodPost,
			endpoint,
			bodyReader,
		)

		if err != nil {
			logs.LogError(err, "Error Preparing Records Request")

			return
		}

		request.Header.Set("Accept", "application/json")
		request.Header.Set("API-VERSION", "1")
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		request.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(request)

		if err != nil {
			logs.LogError(err, "Create Aprimo Record Error")

			return
		}

		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)

		if err != nil {
			logs.LogError(err, "Error Reading Response Body")

			return
		}

		res := string(respBody[:])

		log.Println(res)
	}

	// if err != nil {
	// 	return msgs.SendServerError(err)
	// }

	// endpoint := getEndpointURL("records")

	// resp, err := http.Post(endpoint)

	// return msgs.PrepareResponse(body)
}

func main() {
	lambda.Start(CreateAprimoRecord)
}
