package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/IIP-Design/commons-gateway/utils/data/creds"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
	"github.com/IIP-Design/commons-gateway/utils/queue"
	"github.com/IIP-Design/commons-gateway/utils/security/jwt"
	"github.com/IIP-Design/commons-gateway/utils/turnstile"
)

// verify2FA retrieves that the user provided 2FA
// code matches the values sent to the user.
func verify2FA(id string, code string) bool {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	var storedCode string

	query := "SELECT code FROM mfa WHERE request_id = $1;"
	err = pool.QueryRow(query, id).Scan(&storedCode)

	if err != nil {
		logs.LogError(err, "Retrieve MFA Query Error")
		return false
	}

	return code == storedCode
}

// clear2FA removes the record of a 2FA request after it has been used.
func clear2FA(id string) {
	pool := data.ConnectToDB()
	defer pool.Close()

	query := "DELETE FROM mfa WHERE request_id = $1;"
	_, err := pool.Exec(query, id)

	if err != nil {
		logs.LogError(err, "Delete MFA Request Query Error")
	}
}

// scheduleAccountUnlock initials and SQS message requesting user account unlock in 15 minutes.
func scheduleAccountUnlock(email string) (string, error) {
	var messageId string
	var err error

	event := data.GuestUnlockInitEvent{
		Username: email,
	}

	json, err := json.Marshal(event)

	if err != nil {
		logs.LogError(err, "Failed to Marshal SQS Body")
		return messageId, err
	}

	queueUrl := os.Getenv("UNLOCK_GUEST_ACCOUNT_QUEUE")

	// Send the message to SQS.
	return queue.SendToQueue(string(json), queueUrl)
}

// recordUnsuccessfulLoginAttempt counts the number of failed login attempts
// by a user and locks their account upon a fifth unsuccessful attempt.
func recordUnsuccessfulLoginAttempt(guest string) {
	pool := data.ConnectToDB()
	defer pool.Close()

	currentTime := time.Now()
	var attemptCount int

	query :=
		`UPDATE guests SET login_attempt = login_attempt + 1, login_date = $1
		 WHERE email = $2 RETURNING login_attempt;`
	err := pool.QueryRow(query, currentTime, guest).Scan(&attemptCount)

	if err != nil {
		logs.LogError(err, "Update Login Count Query Error")
	}

	if attemptCount >= 5 {
		query := "UPDATE guests SET locked = true WHERE email = $1"
		_, err := pool.Exec(query, guest)

		if err != nil {
			logs.LogError(err, "Update Lock Status Query Error")
		}

		scheduleAccountUnlock(guest)
	}
}

// handleGrantAccess ensures that a user hash provided a password has matching their
// username and if so, generates a JWT to grant them guest access.
func handleGrantAccess(username string, clientHash string, mfaId string) (msgs.Response, error) {
	if clientHash == "" || username == "" {
		return msgs.Response{StatusCode: 400}, errors.New("data missing from request")
	}

	credentials, err := creds.RetrieveCredentials(username)

	if err != nil {
		return msgs.SendServerError(err)
	}

	if credentials.Hash != clientHash {
		logs.LogError(errors.New("incorrect password"), "Login Error")
		recordUnsuccessfulLoginAttempt(username)
		return msgs.SendCustomError(errors.New("forbidden"), 403)
	} else if credentials.Expired {
		recordUnsuccessfulLoginAttempt(username)
		logs.LogError(errors.New("expired account"), "Login Error")
		return msgs.SendCustomError(errors.New("credentials expired"), 403)
	} else if !credentials.Approved {
		recordUnsuccessfulLoginAttempt(username)
		logs.LogError(errors.New("guest not approved"), "Login Error")
		return msgs.SendCustomError(errors.New("user is not yet approved"), 403)
	} else if credentials.Locked {
		logs.LogError(errors.New("account locked"), "Login Error")
		return msgs.SendCustomError(errors.New("account locked"), 429)
	}

	jwt, err := jwt.FormatJWT(username, credentials.Role, credentials.FirstLogin)

	if err != nil {
		return msgs.SendServerError(err)
	}

	body, err := msgs.MarshalBody(jwt)

	if err != nil {
		return msgs.SendServerError(err)
	}

	// Delete the existing 2FA entry
	clear2FA(mfaId)

	creds.ClearUnsuccessfulLoginAttempts(username)

	return msgs.PrepareResponse(body)
}

// authenticationHandler manages guest user authentication by either generating a JSON web
// token for new authentication sessions or verifying an existing token for an ongoing session.
func authenticationHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	// Retrieve the auth data submitted by the user.
	parsed, err := data.ParseBodyData(event.Body)

	if err != nil {
		return msgs.SendServerError(err)
	}

	clientHash := parsed.Hash
	username := parsed.Username
	mfaId := parsed.MFA.Id
	mfaCode := parsed.MFA.Code

	// Verify that the provided 2FA code is valid.
	verified := verify2FA(mfaId, mfaCode)

	if !verified {
		recordUnsuccessfulLoginAttempt(username)
		logs.LogError(errors.New("submitted 2fa codes does not match"), "Login Error")
		return msgs.SendCustomError(errors.New("forbidden"), 403)
	}

	// Verify the turnstile captcha token
	tokenVerSecretKey := os.Getenv("TOKEN_VERIFICATION_SECRET_KEY")

	if tokenVerSecretKey != "" {
		token := parsed.Token
		remoteIp := event.RequestContext.Identity.SourceIP

		valid, err := turnstile.TokenIsValid(token, remoteIp, tokenVerSecretKey)

		if !valid || err != nil {
			logs.LogError(err, "Turnstile error")
			return msgs.SendServerError(err)
		}
	}

	return handleGrantAccess(username, clientHash, mfaId)
}

func main() {
	lambda.Start(authenticationHandler)
}
