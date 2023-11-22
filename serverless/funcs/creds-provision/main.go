package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/IIP-Design/commons-gateway/utils/data/admins"
	"github.com/IIP-Design/commons-gateway/utils/data/creds"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/email/provision"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
)

// handleInvitation coordinates all the actions associated with inviting a guest user.
func handleInvitation(invite data.Invite) error {
	// Ensure inviter is an active admin user.
	_, adminActive, err := admins.CheckForActiveAdmin(invite.Inviter)

	if err != nil {
		return err
	} else if !adminActive {
		return errors.New("you are not authorized to invite users")
	}

	pass, err := creds.SaveInitialInvite(invite, false)

	if err != nil {
		return err
	}

	_, err = provision.MailProvisionedCreds(invite.Invitee, pass, 0)

	return err
}

// provisionHandler handles the request to grant a guest user temporary credentials. It
// ensures that the required data is present before continuing on to:
//  1. Register the invitation
//  2. Provision credentials for the guest user
//  3. Initiate the admin and guest user notifications
func provisionHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	invite, err := data.ExtractInvite(event.Body)

	if err != nil {
		return msgs.SendServerError(err)
	}

	err = handleInvitation(invite)

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(provisionHandler)
}
