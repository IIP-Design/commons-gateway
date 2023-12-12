package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/IIP-Design/commons-gateway/utils/data/admins"
	"github.com/IIP-Design/commons-gateway/utils/data/creds"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/email/provision"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
)

// handleInvitation coordinates all the actions associated with inviting a guest user.
func handleInvitation(invite data.Invite) error {
	// Ensure inviter is an active admin user.
	_, adminActive, err := admins.CheckForActiveAdmin(invite.Inviter)

	if err != nil {
		logs.LogError(err, "Admin Check Error")
		return err
	} else if !adminActive {
		err = fmt.Errorf("the user %s is not authorized to add users", invite.Inviter)

		logs.LogError(err, "Admin Check Error")
		return err
	}

	fmt.Printf("Registering the invitation of %s by %s\n", invite.Invitee.Email, invite.Inviter)

	pass, err := creds.SaveInitialInvite(invite, false)

	if err != nil {
		logs.LogError(err, "Save Credentials Error")
		return err
	}

	fmt.Printf("Sending %s their temporary credentials\n", invite.Invitee.Email)

	_, err = provision.MailProvisionedCreds(invite.Invitee, pass, 0)

	if err != nil {
		logs.LogError(err, "Mail Credentials Error")
	}

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
		logs.LogError(err, "Extract Invite Error")
		return msgs.SendServerError(err)
	}

	err = handleInvitation(invite)

	if err != nil {
		logs.LogError(err, "Handle Invite Error")

		if err.Error() == "user already exists" {
			return msgs.SendCustomError(err, 409)
		}

		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(provisionHandler)
}
