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

	err := testHelpers.SetUpTestDb()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	exitVal := m.Run()

	testHelpers.TearDownTestDb()
	testHelpers.CleanupInvites(testHelpers.ExampleGuest2["email"])

	os.Exit(exitVal)
}

func TestGoodData(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: makeJsonBody(testHelpers.ExampleGuest2["email"], testHelpers.ExampleGuest["email"]),
	}

	resp, err := proposalHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("proposalHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	pending, err := testHelpers.CheckGuestPending(testHelpers.ExampleGuest2["email"])
	if !pending || err != nil {
		t.Fatalf("CheckGuestPending result %t/%v, want true/nil", pending, err)
	}
}

func TestBadProposer(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: makeJsonBody(testHelpers.ExampleGuest2["email"], ""),
	}

	resp, err := proposalHandler(context.TODO(), event)
	if resp.StatusCode != 500 || err != nil {
		t.Fatalf("proposalHandler result %d/%v, want 500/nil", resp.StatusCode, err)
	}
}

func TestBadInvite(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: makeJsonBody("", testHelpers.ExampleGuest["email"]),
	}

	resp, err := proposalHandler(context.TODO(), event)
	if resp.StatusCode != 500 || err != nil {
		t.Fatalf("proposalHandler result %d/%v, want 500/nil", resp.StatusCode, err)
	}
}

func TestProposerMiss(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: makeJsonBody(testHelpers.ExampleGuest2["email"], "fake@test.fail"),
	}

	resp, err := proposalHandler(context.TODO(), event)
	if resp.StatusCode != 500 || err != nil {
		t.Fatalf("proposalHandler result %d/%v, want 500/nil", resp.StatusCode, err)
	}
}

func makeJsonBody(inviteeEmail string, proposerEmail string) string {
	return fmt.Sprintf(`{
		"invitee": {
			"email": "%s",
			"givenName": "%s",
			"familyName": "%s",
			"role": "guest",
			"team": "%s"
		},
		"proposer": "%s",
		"expiration": "%s"
	}`,
		inviteeEmail,
		testHelpers.ExampleGuest2["first_name"],
		testHelpers.ExampleGuest2["last_name"],
		testHelpers.ExampleTeam["id"],
		proposerEmail,
		testHelpers.FarFutureDateStr())
}
