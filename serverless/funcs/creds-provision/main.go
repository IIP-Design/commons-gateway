package main

import (
	"context"
	"errors"
	"fmt"

	data "github.com/IIP-Design/commons-gateway/utils/data"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// handleInvitation coordinates all the actions associated with inviting a guest user.
func handleInvitation(invite data.Invite) error {
	var err error

	// Ensure inviter is an active admin user.
	adminActive, err := data.CheckForActiveAdmin(invite.Inviter)

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

	err = data.SaveCredentials(invite.Invitee, hash, salt)

	if err != nil {
		return errors.New("something went wrong - credential generation failed")
	}

	// Record the invitation - has to follow cred generation due to foreign key constraint
	err = data.SaveInvite(invite.Inviter, invite.Invitee.Email)

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
	var msg string

	invite, err := data.ExtractInvite(event.Body)

	if err != nil {
		return msgs.Response{StatusCode: 500}, err
	}

	err = handleInvitation(invite)

	if err != nil {
		return msgs.Response{StatusCode: 500}, err
	} else {
		msg = "success"
	}

	return msgs.PrepareResponse(msg)
}

func main() {
	lambda.Start(ProvisionHandler)
}
