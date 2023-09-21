package main

import (
	"context"
	"errors"
	"os"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/guests"
	"github.com/IIP-Design/commons-gateway/utils/email/provision"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
	"github.com/IIP-Design/commons-gateway/utils/security/hashing"
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

	// TODO - email URL
	sourceEmail := os.Getenv("SOURCE_EMAIL_ADDRESS")
	redirectUrl := os.Getenv("EMAIL_REDIRECT_URL")
	err = provision.MailProvisionedCreds(sourceEmail, provision.ProvisionCredsData{
		Invitee:     invitee,
		TmpPassword: pass,
		Url:         redirectUrl,
	})
	if err != nil {
		return msgs.SendServerError(err)
	}

	return msgs.SendSuccessMessage()
}

func main() {
	lambda.Start(GuestAcceptHandler)
}
