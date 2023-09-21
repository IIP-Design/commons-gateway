package main

import (
	"context"
	"errors"
	"os"

	"github.com/IIP-Design/commons-gateway/utils/data/admins"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/invites"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
	"github.com/IIP-Design/commons-gateway/utils/security/hashing"
	"github.com/IIP-Design/commons-gateway/utils/security/jwt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// handleInvitation coordinates all the actions associated with inviting a guest user.
func handleProposedInvitation(invite data.Invite) error {
	var err error

	// Ensure proposer is an active admin user.
	proposer, active, err := admins.CheckForGuestAdmin(invite.Proposer)

	if err != nil {
		return err
	} else if !active {
		return errors.New("you are not authorized to propose user invites")
	}

	// Ensure invitee doesn't already have access.
	_, guestHasAccess, err := data.CheckForExistingUser(invite.Invitee.Email, "guests")

	if err != nil {
		return err
	} else if guestHasAccess {
		return errors.New("guest user already has access")
	}

	// Generate credentials
	pass, salt := hashing.GenerateCredentials()
	hash := hashing.GenerateHash(pass, salt)

	// PASSWORD IS UNRECOVERABLE
	err = invites.SaveCredentials(invite.Invitee, invite.Expires, hash, salt)

	if err != nil {
		return errors.New("something went wrong - credential generation failed")
	}

	// Record the invitation - has to follow cred generation due to foreign key constraint
	err = invites.SaveInvite(invite.Proposer, invite.Invitee.Email, true)

	if err != nil {
		return errors.New("something went wrong - saving invite failed")
	}

	// TODO - email URL
	sourceEmail := os.Getenv("SOURCE_EMAIL_ADDRESS")
	err = MailProposedCreds(sourceEmail, RequestSupportStaffData{
		Invitee:  invite.Invitee,
		Proposer: proposer,
		Url:      "/invites",
	})

	return err
}

// ProvisionHandler handles the request to grant a guest user temporary credentials. It
// ensures that the required data is present before continuing on to:
//  1. Register the proposed invitation
//  2. Provision preliminary credentials for the guest user
//  3. Initiate the admin and guest user notifications
func ProposalHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	code, err := jwt.RequestIsAuthorized(event, []string{"guest"})
	if err != nil {
		return msgs.SendAuthError(err, code)
	}

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
