package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/IIP-Design/commons-gateway/utils/data/admins"
	"github.com/IIP-Design/commons-gateway/utils/data/creds"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/email/propose"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// handleInvitation coordinates all the actions associated with inviting a guest user.
func handleProposedInvitation(invite data.Invite) error {
	var err error

	// Ensure proposer is an active admin user.
	proposer, active, err := admins.CheckForGuestAdmin(invite.Proposer)

	if err != nil {
		logs.LogError(err, "Admin Check Error")
		return err
	} else if !active {
		err = fmt.Errorf("the user %s is not authorized propose user invites", proposer.Email)

		logs.LogError(err, "Admin Check Error")
		return errors.New("you are not authorized to propose user invites")
	}

	fmt.Printf("Registering the invitation of %s by %s\n", invite.Invitee.Email, proposer.Email)

	_, err = creds.SaveInitialInvite(invite, true)

	if err != nil {
		logs.LogError(err, "Save Credentials Error")
		return err
	}

	fmt.Printf("Sending %s their temporary credentials\n", invite.Invitee.Email)

	err = propose.MailProposedInvite(proposer, invite.Invitee)

	return err
}

// ProvisionHandler handles the request to grant a guest user temporary credentials. It
// ensures that the required data is present before continuing on to:
//  1. Register the proposed invitation
//  2. Provision preliminary credentials for the guest user
//  3. Initiate the admin and guest user notifications
func proposalHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	invite, err := data.ExtractInvite(event.Body)

	if err != nil {
		logs.LogError(err, "Extract Invite Error")
		return msgs.SendServerError(err)
	}

	err = handleProposedInvitation(invite)

	if err != nil {
		logs.LogError(err, "Handle Proposed Invite Error")
		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(proposalHandler)
}
