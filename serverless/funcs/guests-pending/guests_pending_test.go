package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	testConfig "github.com/IIP-Design/commons-gateway/test/config"
	testHelpers "github.com/IIP-Design/commons-gateway/test/helpers"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/aws/aws-lambda-go/events"
)

type DataBody struct {
	Data []map[string]string `json:"data"`
}

func TestMain(m *testing.M) {
	testConfig.ConfigureDb()

	testHelpers.TearDownTestDb()
	err := testHelpers.SetUpTestDb()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	exitVal := m.Run()

	cleanupInvites()
	testHelpers.TearDownTestDb()

	os.Exit(exitVal)
}

func TestGetGuestsWithTeamNoData(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"team":"%s"}`, testHelpers.ExampleTeam["id"]),
	}

	resp, err := getPendingInvitesHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("getPendingInvitesHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	body := resp.Body
	guests, err := deserializeBody(body)

	if err != nil || len(guests) != 0 {
		t.Fatalf("Data is ill-formed: %v/%d", err, len(guests))
	}
}

func TestGetGuestsNoTeamNoData(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"team":"%s"}`, ""),
	}

	resp, err := getPendingInvitesHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("getPendingInvitesHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	body := resp.Body
	guests, err := deserializeBody(body)

	if err != nil || len(guests) != 0 {
		t.Fatalf("Data is ill-formed: %v/%d", err, len(guests))
	}
}

func TestGetGuestsWithTeamWithData(t *testing.T) {
	testHelpers.AddPendingGuest()
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"team":"%s"}`, testHelpers.ExampleTeam["id"]),
	}

	resp, err := getPendingInvitesHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("getPendingInvitesHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	body := resp.Body
	guests, err := deserializeBody(body)

	if err != nil || len(guests) == 0 || guests[0]["email"] != testHelpers.ExampleGuest2["email"] {
		t.Fatalf("Data is ill-formed: %v/%d", err, len(guests))
	}
}

func TestGetGuestsNoTeamWithData(t *testing.T) {
	err := testHelpers.AddPendingGuest()
	if err != nil {
		t.Fatalf(`Invite error: %v`, err)
	}

	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"team":"%s"}`, ""),
	}

	resp, err := getPendingInvitesHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("getPendingInvitesHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	body := resp.Body
	guests, err := deserializeBody(body)

	if err != nil || len(guests) == 0 || guests[0]["email"] != testHelpers.ExampleGuest2["email"] {
		t.Fatalf("Data is ill-formed: %v/%d", err, len(guests))
	}
}

func deserializeBody(body string) ([]map[string]string, error) {
	var parsed DataBody

	b := []byte(body)
	err := json.Unmarshal(b, &parsed)

	return parsed.Data, err
}

func cleanupInvites() {
	pool := data.ConnectToDB()
	defer pool.Close()

	pool.Exec("DELETE FROM invites WHERE invitee = $1", testHelpers.ExampleGuest2["email"])
	pool.Exec("DELETE FROM guests WHERE email = $1", testHelpers.ExampleGuest2["email"])
}
