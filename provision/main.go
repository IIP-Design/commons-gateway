package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// EventData describes the data that the Lambda function expects to receive.
type EventData struct {
	Invitee string `json:"invitee"`
	Inviter string `json:"inviter"`
}

// handleInvitation coordinates all the actions associated with inviting a guest user.
func handleInvitation(adminEmail string, guestEmail string) error {
	var err error

	guestHasAccess, err := checkForExistingAccess(guestEmail)

	if err != nil {
		return err
	}

	if guestHasAccess {
		return errors.New("guest user already has access")
	} else {
		// Record the invitation
		err = saveInvite(adminEmail, guestEmail)

		if err != nil {
			return errors.New("something went wrong - saving invite failed")
		}

		// Generate credentials
		pass, salt := generateOTP()
		hash := generateHash(pass, salt)

		err = saveCredentials(guestEmail, hash, salt)

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
func ProvisionHandler(ctx context.Context, data EventData) (Response, error) {
	var msg string

	inviter := data.Inviter
	invitee := data.Invitee

	if inviter == "" || invitee == "" {
		return Response{StatusCode: 400}, errors.New("data missing from request")
	}

	err := handleInvitation(inviter, invitee)

	if err != nil {
		return Response{StatusCode: 500}, err
	} else {
		msg = "success"
	}

	return prepareResponse(msg)
}

func main() {
	lambda.Start(ProvisionHandler)
}
