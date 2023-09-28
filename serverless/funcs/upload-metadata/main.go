package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
	"github.com/IIP-Design/commons-gateway/utils/security/jwt"

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

// retrieveUserId converts an admin or guest id value into a user id value.
func retrieveUserId(email string, pool *sql.DB) (string, error) {
	var err error
	var userId string

	// Check if user exists in the admins table
	_, isAdmin, err := data.CheckForExistingUser(email, "admins")

	if err == nil && isAdmin {
		// Retrieve admin's the user_id from the all_users table
		err := pool.QueryRow(`SELECT user_id FROM all_users WHERE admin_id = $1`, email).Scan(&userId)

		if err != nil {
			logs.LogError(err, "Select Admin's User Id Query Error")
		}

		return userId, err
	} else if err == nil && !isAdmin {
		// If not found in admin table, look in the guests table.
		_, isGuest, err := data.CheckForExistingUser(email, "guests")

		if err != nil {
			logs.LogError(err, "Check for Guest Query Error")

			return userId, err
		}

		if isGuest {
			// Retrieve the new user id from the all_users table
			err := pool.QueryRow(`SELECT user_id FROM all_users WHERE guest_id = $1`, email).Scan(&userId)

			if err != nil {
				logs.LogError(err, "Select Guest's User Id Query Error")
			}
		} else {
			err = errors.New("user is neither an admin nor a guest")
		}

		return userId, err
	} else if err != nil {
		logs.LogError(err, "Check for Admin User Query Error")
	}

	return userId, err
}

// createUploadRecord opens a connection to the database and add a new upload record.
func createUploadRecord(s3Id string, user string, teamId string, fileType string, description string) error {
	pool := data.ConnectToDB()
	defer pool.Close()

	currentTime := time.Now()

	id, err := retrieveUserId(user, pool)

	if err != nil {
		logs.LogError(err, "Retrieve User Id Query Error")

		return err
	}

	query := "INSERT INTO uploads ( s3_id, user_id, team_id, file_type, description, date_uploaded ) VALUES ( $1, $2, $3, $4, $5, $6 )"
	_, err = pool.Exec(query, s3Id, id, teamId, fileType, description, currentTime)

	if err != nil {
		logs.LogError(err, "Create Upload Record Query Error")
	}

	return err
}

// NewTeamHandler handles the request to add a new team for uploading. It
// ensures that the required data is present before continuing on to recording
// the team name and setting it to active.
func NewUploadHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	code, err := jwt.RequestIsAuthorized(event, []string{"super admin", "admin", "guest admin", "guest"})

	if err != nil {
		return msgs.SendAuthError(err, code)
	}

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
