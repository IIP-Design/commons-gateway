package main

import (
	"context"
	"errors"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/guests"
	"github.com/IIP-Design/commons-gateway/utils/email/propose"
	"github.com/IIP-Design/commons-gateway/utils/email/provision"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
	"github.com/IIP-Design/commons-gateway/utils/security/jwt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func GuestReauthHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	// Need client role to determine reauth logic
	scope, err := jwt.ExtractClientRole(event.Headers["authorization"])
	if err != nil {
		return msgs.SendServerError(err)
	}

	clientIsGuestAdmin := (scope == "guest admin")

	guest, err := data.ExtractReauth(event.Body)
	if err != nil {
		return msgs.SendServerError(err)
	}

	// Ensure that the user we intend to modify exists.
	user, userExists, err := data.CheckForExistingUser(guest.Email, "guests")

	if err != nil {
		return msgs.SendServerError(err)
	} else if !userExists {
		return msgs.SendServerError(errors.New("this user has not been registered"))
	}

	// Try to reauthorize
	pass, status, err := guests.Reauthorize(guest, clientIsGuestAdmin)

	// May indicate a conflict (they have a pending request) or server error
	if err != nil {
		return msgs.SendCustomError(err, status)
	}

	// For guest admins, we always need to email an admin to approve the new creds
	if clientIsGuestAdmin {
		proposer, _, err := data.CheckForExistingUser(guest.Admin, "guests")
		if err != nil {
			return msgs.SendServerError(err)
		}

		err = propose.MailProposedCreds(user, proposer)
		if err != nil {
			return msgs.SendServerError(err)
		}
	} else if pass != "" {
		// For admins, only send an email if they need to re-up their password
		err = provision.MailProvisionedCreds(user, pass)
		if err != nil {
			return msgs.SendServerError(err)
		}
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(GuestReauthHandler)
}
