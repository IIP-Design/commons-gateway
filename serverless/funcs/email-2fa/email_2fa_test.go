package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	testConfig "github.com/IIP-Design/commons-gateway/test/config"
	testHelpers "github.com/IIP-Design/commons-gateway/test/helpers"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/aws/aws-lambda-go/events"
)

const (
	CODE       = "123456"
	MESSAGE_ID = "98765"
)

func TestMain(m *testing.M) {
	testConfig.ConfigureEmail()

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestFmtEmailBody(t *testing.T) {
	pattern := regexp.MustCompile(fmt.Sprintf(`(?m)<p>%s</p>`, CODE))
	body := formatEmailBody(makeUser(), CODE)

	match := pattern.MatchString(body)

	if !match {
		t.Fatal("Failed to match email body")
	}
}

func TestFmtEmail(t *testing.T) {
	sourceEmail := os.Getenv("SOURCE_EMAIL_ADDRESS")
	user := makeUser()
	email := formatEmail(user, CODE, sourceEmail)

	if email.Destination.ToAddresses[0] != user.Email {
		t.Fatalf(`Email ill formed: %s, expected %s`, email.Destination.ToAddresses[0], user.Email)
	}
}

func TestSendEmail(t *testing.T) {
	eventBody := fmt.Sprintf(`{"code":"%s","user":{"email":"%s","givenName":"%s","familyName":"%s"}}`,
		CODE, testHelpers.ExampleGuest["email"], testHelpers.ExampleGuest["first_name"], testHelpers.ExampleGuest["last_name"])
	event := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId: MESSAGE_ID,
				Body:      eventBody,
			},
		},
	}

	err := email2FAHandler(context.TODO(), event)
	if err != nil {
		t.Fatalf(`Error: %v, want nil`, err)
	}
}

func TestSendEmailBadData(t *testing.T) {
	eventBody := fmt.Sprintf(`{code:"user":{"email":"%s","givenName":"%s","familyName":"%s"}}`,
		testHelpers.ExampleGuest["email"], testHelpers.ExampleGuest["first_name"], testHelpers.ExampleGuest["last_name"])
	event := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId: MESSAGE_ID,
				Body:      eventBody,
			},
		},
	}

	err := email2FAHandler(context.TODO(), event)
	if err == nil {
		t.Fatal(`Error: nil, want not nil`)
	}
}

func makeUser() data.User {
	return data.User{
		Email:     testHelpers.ExampleGuest["email"],
		NameFirst: testHelpers.ExampleGuest["first_name"],
		NameLast:  testHelpers.ExampleGuest["last_name"],
	}
}
