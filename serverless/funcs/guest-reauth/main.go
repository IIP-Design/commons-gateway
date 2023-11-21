package main

import (
	"context"
	"errors"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/guests"
	"github.com/IIP-Design/commons-gateway/utils/data/users"
	"github.com/IIP-Design/commons-gateway/utils/email/propose"
	"github.com/IIP-Design/commons-gateway/utils/email/provision"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
	"github.com/IIP-Design/commons-gateway/utils/security/jwt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func guestReauthHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	// Need client role to determine reauthorization logic
	scope, err := jwt.ExtractClientRole(event.Headers["Authorization"])
	if err != nil {
		return msgs.SendServerError(err)
	}

	clientIsGuestAdmin := (scope == "guest admin")

	guest, err := data.ExtractReauth(event.Body)

	if err != nil {
		logs.LogError(err, "Data Extraction Error")

		return msgs.SendServerError(err)
	}

	// Ensure that the user we intend to modify exists.
	user, userExists, err := users.CheckForExistingUser(guest.Email, "guests")

	if err != nil {
		logs.LogError(err, "User Check Error")
		return msgs.SendServerError(err)
	} else if !userExists {
		logs.LogError(err, "User Not Found Error")
		return msgs.SendCustomError(errors.New("this user has not been registered"), 404)
	}

	// Try to reauthorize
	pass, status, err := guests.Reauthorize(guest, clientIsGuestAdmin)

	// May indicate a conflict (they have a pending request) or server error
	if err != nil {
		logs.LogError(err, "User Reauthorization Error")
		return msgs.SendCustomError(err, status)
	} else if status >= 400 {
		err := errors.New("user reauthorization conflict")
		return msgs.SendCustomError(err, status)
	}

	// For guest admins, we always need to email an admin to approve the new creds
	if clientIsGuestAdmin {
		proposer, _, err := users.CheckForExistingUser(guest.Admin, "guests")
		if err != nil {
			return msgs.SendServerError(err)
		}

		err = propose.MailProposedCreds(user, proposer)

		if err != nil {
			return msgs.SendServerError(err)
		}
	} else if pass != "" {
		// For admins, only send an email if they need to re-up their password
		_, err = provision.MailProvisionedCreds(user, pass, 1)

		if err != nil {
			return msgs.SendServerError(err)
		}
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(guestReauthHandler)
}
