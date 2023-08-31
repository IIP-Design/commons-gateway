package main

import (
	"context"
	"encoding/json"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RequestBody struct {
	S3Id        string `json:"key"`
	User        string `json:"email"`
	TeamId      string `json:"team"`
	FileType    string `json:"fileType"`
	Description string `json:"description"`
}

func parseRequest(body string) (RequestBody, error) {
	var parsed RequestBody

	b := []byte(body)
	err := json.Unmarshal(b, &parsed)

	if err != nil {
		logs.LogError(err, "Failed to Unmarshal Body")
	}

	return parsed, err
}

func createUploadRecord(s3Id string, user string, teamId string, fileType string, description string) error {
	pool := data.ConnectToDB()
	defer pool.Close()

	query := "INSERT INTO uploads ( s3_id, user_id, team_id, file_type, description ) VALUES ( $1, $2, $3, $4, $5 )"
	_, err := pool.Exec(query, s3Id, user, teamId, fileType, description)

	if err != nil {
		logs.LogError(err, "Create Upload Record Query Error")
	}

	return err
}

// NewTeamHandler handles the request to add a new team for uploading. It
// ensures that the required data is present before continuing on to recording
// the team name and setting it to active.
func NewUploadHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	parsed, err := parseRequest(event.Body)

	s3Id := parsed.S3Id
	user := parsed.User
	teamId := parsed.TeamId
	fileType := parsed.FileType
	description := parsed.Description

	if err != nil {
		return msgs.SendServerError(err)
	}

	err = createUploadRecord(s3Id, user, teamId, fileType, description)

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.PrepareResponse([]byte("success"))
}

func main() {
	lambda.Start(NewUploadHandler)
}
