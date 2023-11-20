package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	testConfig "github.com/IIP-Design/commons-gateway/test/config"
	testHelpers "github.com/IIP-Design/commons-gateway/test/helpers"
	"github.com/aws/aws-lambda-go/events"
)

func TestMain(m *testing.M) {
	testConfig.ConfigureDb()
	testConfig.ConfigureEmail()

	testHelpers.TearDownTestDb()
	err := testHelpers.SetUpTestDb()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = testHelpers.AddPendingGuest()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	exitVal := m.Run()

	testHelpers.CleanupInvites(testHelpers.ExampleGuest2["email"])
	testHelpers.TearDownTestDb()

	os.Exit(exitVal)
}

func TestApprove(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: makeJsonBody(testHelpers.ExampleGuest2["email"], testHelpers.ExampleAdmin["email"]),
	}

	resp, err := guestAcceptHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("guestAcceptHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	pending, err := testHelpers.CheckGuestPending(testHelpers.ExampleGuest2["email"])
	if pending || err != nil {
		t.Fatalf("CheckGuestPending result %t/%v, want true/nil", pending, err)
	}
}

func TestUserMiss(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: makeJsonBody("fake@test.fail", testHelpers.ExampleAdmin["email"]),
	}

	resp, err := guestAcceptHandler(context.TODO(), event)
	if resp.StatusCode != 404 || err != nil {
		t.Fatalf("guestAcceptHandler result %d/%v, want 404/nil", resp.StatusCode, err)
	}
}

func makeJsonBody(inviteeEmail string, inviterEmail string) string {
	return fmt.Sprintf(`{
		"inviteeEmail": "%s",
		"inviterEmail": "%s"
	}`,
		inviteeEmail,
		inviterEmail)
}
