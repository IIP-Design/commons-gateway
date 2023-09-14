package main

import (
	"context"
	"errors"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/guests"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
	"github.com/IIP-Design/commons-gateway/utils/security/jwt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// GuestAcceptHandler accepts a request to invite an external partner.
func GuestAcceptHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	code, err := jwt.RequestIsAuthorized(event, []string{"super admin", "admin"})
	if err != nil {
		return msgs.SendAuthError(err, code)
	}

	guest, err := data.ExtractAcceptInvite(event.Body)

	if err != nil {
		return msgs.SendServerError(err)
	}

	// Ensure that the user we intend to modify exists.
	userExists, err := data.CheckForExistingUser(guest.Invitee, "guests")

	if err != nil {
		return msgs.SendServerError(err)
	} else if !userExists {
		return msgs.SendServerError(errors.New("this user has not been invited"))
	}

	err = guests.AcceptGuest(guest)

	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(GuestAcceptHandler)
}
