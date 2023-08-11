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
func handleInvitation(adminEmail string, guestEmail string) error {
	var err error

	// Ensure inviter is an active admin user.
	adminActive, err := data.CheckForActiveAdmin(adminEmail)

	if err != nil {
		return err
	} else if !adminActive {
		return errors.New("you are not authorized to invite users")
	}

	// Ensure invitee doesn't already have access.
	guestHasAccess, err := data.CheckForExistingUser(guestEmail, "credentials")

	if err != nil {
		return err
	} else if guestHasAccess {
		return errors.New("guest user already has access")
	} else {
		// Record the invitation
		err = data.SaveInvite(adminEmail, guestEmail)

		if err != nil {
			return errors.New("something went wrong - saving invite failed")
		}

		// Generate credentials
		pass, salt := generateCredentials()
		hash := generateHash(pass, salt)

		err = data.SaveCredentials(guestEmail, hash, salt)

		if err == nil {
			// TODO - send password
			fmt.Printf("Your password is %s", pass)
		} else {
			return errors.New("something went wrong - credential generation failed")
		}
	}

	return err
}

// ProvisionHandler handles the request to grant a guest user temporary credentials. It
// ensures that the required data is present before continuing on to:
//  1. Register the invitation
//  2. Provision credentials for the guest user
//  3. Initiate the admin and guest user notifications
func ProvisionHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	var msg string

	parsed, err := data.ParseBodyData(event.Body)

	inviter := parsed.Inviter
	invitee := parsed.Invitee

	if err != nil {
		return msgs.Response{StatusCode: 500}, err
	} else if inviter == "" || invitee == "" {
		return msgs.Response{StatusCode: 400}, errors.New("data missing from request")
	}

	err = handleInvitation(inviter, invitee)

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
