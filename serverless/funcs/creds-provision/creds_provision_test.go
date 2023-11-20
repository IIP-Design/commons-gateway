package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	testConfig "github.com/IIP-Design/commons-gateway/test/config"
	testHelpers "github.com/IIP-Design/commons-gateway/test/helpers"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
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
	cleanupInvites()

	os.Exit(exitVal)
}

func TestGoodData(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: makeJsonBody(testHelpers.ExampleGuest2["email"], testHelpers.ExampleAdmin["email"]),
	}

	resp, err := provisionHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("provisionHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	pending, err := checkGuestPending(testHelpers.ExampleGuest2["email"])
	if pending || err != nil {
		t.Fatalf("checkGuestPending result %t/%v, want false/nil", pending, err)
	}
}

func TestBadAdmin(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: makeJsonBody(testHelpers.ExampleGuest2["email"], ""),
	}

	resp, err := provisionHandler(context.TODO(), event)
	if resp.StatusCode != 500 || err != nil {
		t.Fatalf("provisionHandler result %d/%v, want 500/nil", resp.StatusCode, err)
	}
}

func TestBadInvite(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: makeJsonBody("", testHelpers.ExampleAdmin["email"]),
	}

	resp, err := provisionHandler(context.TODO(), event)
	if resp.StatusCode != 500 || err != nil {
		t.Fatalf("provisionHandler result %d/%v, want 500/nil", resp.StatusCode, err)
	}
}

func TestAdminMiss(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: makeJsonBody(testHelpers.ExampleGuest2["email"], "fake@test.fail"),
	}

	resp, err := provisionHandler(context.TODO(), event)
	if resp.StatusCode != 500 || err != nil {
		t.Fatalf("provisionHandler result %d/%v, want 500/nil", resp.StatusCode, err)
	}
}

func makeJsonBody(inviteeEmail string, adminEmail string) string {
	return fmt.Sprintf(`{
		"invitee": {
			"email": "%s",
			"givenName": "%s",
			"familyName": "%s",
			"role": "guest",
			"team": "%s"
		},
		"inviter": "%s",
		"expiration": "2025-12-01T12:00:00Z"
	}`,
		inviteeEmail,
		testHelpers.ExampleGuest2["first_name"],
		testHelpers.ExampleGuest2["last_name"],
		testHelpers.ExampleTeam["id"],
		adminEmail)
}

func checkGuestPending(email string) (bool, error) {
	pool := data.ConnectToDB()
	defer pool.Close()

	var pending bool
	query := `SELECT pending FROM invites WHERE invitee = $1`
	err := pool.QueryRow(query, email).Scan(&pending)

	return pending, err
}

func cleanupInvites() {
	pool := data.ConnectToDB()
	defer pool.Close()

	pool.Exec("DELETE FROM invites WHERE invitee = $1", testHelpers.ExampleGuest2["email"])
	pool.Exec("DELETE FROM guests WHERE email = $1", testHelpers.ExampleGuest2["email"])
}
