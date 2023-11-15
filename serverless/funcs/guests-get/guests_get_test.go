package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	testConfig "github.com/IIP-Design/commons-gateway/test/config"
	testHelpers "github.com/IIP-Design/commons-gateway/test/helpers"
	"github.com/IIP-Design/commons-gateway/utils/types"
	"github.com/aws/aws-lambda-go/events"
)

type DataBody struct {
	Data []types.GuestUser `json:"data"`
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

	testHelpers.TearDownTestDb()

	os.Exit(exitVal)
}

func TestGetGuestsNoTeam(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"role":"%s"}`,
			testHelpers.ExampleGuest["role"]),
	}

	resp, err := getGuestsHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("getGuestsHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	body := resp.Body
	guests, err := deserializeBody(body)

	if err != nil || len(guests) == 0 {
		t.Fatalf("Data is ill-formed: %v/%d", err, len(guests))
	}
}

func TestGetGuestsRealTeam(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"role":"%s","team":"%s"}`,
			testHelpers.ExampleGuest["role"], testHelpers.ExampleTeam["id"]),
	}

	resp, err := getGuestsHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("getGuestsHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	body := resp.Body
	guests, err := deserializeBody(body)

	if err != nil || len(guests) == 0 {
		t.Fatalf("Data is ill-formed: %v/%d", err, len(guests))
	}
}

func TestGetGuestsFakeTeam(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"role":"%s","team":"%s"}`,
			testHelpers.ExampleGuest["role"], "ERROR"),
	}

	resp, err := getGuestsHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("getGuestsHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	body := resp.Body
	guests, err := deserializeBody(body)

	if err != nil || len(guests) != 0 {
		t.Fatalf("Data is ill-formed: %v/%d", err, len(guests))
	}
}

func deserializeBody(body string) ([]types.GuestUser, error) {
	var parsed DataBody

	b := []byte(body)
	err := json.Unmarshal(b, &parsed)

	return parsed.Data, err
}
