package main

import (
	"context"
	"errors"

	"github.com/IIP-Design/commons-gateway/utils/data/admins"
	"github.com/IIP-Design/commons-gateway/utils/data/creds"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/email/propose"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
	"github.com/IIP-Design/commons-gateway/utils/types"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// handleInvitation coordinates all the actions associated with inviting a guest user.
func handleProposedInvitation(invite types.Invite) error {
	var err error

	// Ensure proposer is an active admin user.
	proposer, active, err := admins.CheckForGuestAdmin(invite.Proposer)

	if err != nil {
		return err
	} else if !active {
		return errors.New("you are not authorized to propose user invites")
	}

	_, err = creds.SaveInitialInvite(invite, true)
	if err != nil {
		return err
	}

	err = propose.MailProposedCreds(invite.Invitee, proposer)

	return err
}

// ProvisionHandler handles the request to grant a guest user temporary credentials. It
// ensures that the required data is present before continuing on to:
//  1. Register the proposed invitation
//  2. Provision preliminary credentials for the guest user
//  3. Initiate the admin and guest user notifications
func ProposalHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	invite, err := data.ExtractInvite(event.Body)

	if err != nil {
		return msgs.SendServerError(err)
	}

	err = handleProposedInvitation(invite)

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(ProposalHandler)
}
