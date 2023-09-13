package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/IIP-Design/commons-gateway/utils/data/admins"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/invites"
	"github.com/IIP-Design/commons-gateway/utils/jwt"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// handleInvitation coordinates all the actions associated with inviting a guest user.
func handleInvitation(invite data.Invite) error {
	var err error

	// Ensure inviter is an active admin user.
	adminActive, err := admins.CheckForActiveAdmin(invite.Inviter)

	if err != nil {
		return err
	} else if !adminActive {
		return errors.New("you are not authorized to invite users")
	}

	// Ensure invitee doesn't already have access.
	guestHasAccess, err := data.CheckForExistingUser(invite.Invitee.Email, "guests")

	if err != nil {
		return err
	} else if guestHasAccess {
		return errors.New("guest user already has access")
	}

	// Generate credentials
	pass, salt := generateCredentials()
	hash := generateHash(pass, salt)

	err = invites.SaveCredentials(invite.Invitee, invite.Expires, hash, salt)

	if err != nil {
		return errors.New("something went wrong - credential generation failed")
	}

	// Record the invitation - has to follow cred generation due to foreign key constraint
	err = invites.SaveInvite(invite.Inviter, invite.Invitee.Email)

	if err != nil {
		return errors.New("something went wrong - saving invite failed")
	}

	// TODO - email password
	fmt.Printf("Your password is %s", pass)
	return err
}

// ProvisionHandler handles the request to grant a guest user temporary credentials. It
// ensures that the required data is present before continuing on to:
//  1. Register the invitation
//  2. Provision credentials for the guest user
//  3. Initiate the admin and guest user notifications
func ProvisionHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	code, err := jwt.RequestIsAuthorized(event, []string{"super admin", "admin"})
	if err != nil {
		return msgs.SendAuthError(err, code)
	}

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
	lambda.Start(ProvisionHandler)
}
