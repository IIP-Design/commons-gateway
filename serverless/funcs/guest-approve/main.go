package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/guests"
	"github.com/IIP-Design/commons-gateway/utils/email/provision"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
	"github.com/IIP-Design/commons-gateway/utils/security/hashing"
)

// GuestAcceptHandler accepts a request to invite an external partner.
func GuestAcceptHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	guest, err := data.ExtractAcceptInvite(event.Body)

	if err != nil {
		return msgs.SendServerError(err)
	}

	// Ensure that the user we intend to modify exists.
	invitee, userExists, err := data.CheckForExistingUser(guest.Invitee, "guests")

	if err != nil {
		return msgs.SendServerError(err)
	} else if !userExists {
		return msgs.SendServerError(errors.New("this user has not been invited"))
	}

	// Regenerate credentials
	pass, salt := hashing.GenerateCredentials()
	hash := hashing.GenerateHash(pass, salt)

	err = guests.AcceptGuest(guest, hash, salt)

	if err != nil {
		return msgs.SendServerError(err)
	}

	err = provision.MailProvisionedCreds(invitee, pass, 1)

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(GuestAcceptHandler)
}
