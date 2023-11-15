package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	testConfig "github.com/IIP-Design/commons-gateway/test/config"
	testHelpers "github.com/IIP-Design/commons-gateway/test/helpers"
	"github.com/aws/aws-lambda-go/events"
)

type SaltData struct {
	Salt      string   `json:"salt"`
	PrevSalts []string `json:"prevSalts"`
}

type DataBody struct {
	Data SaltData `json:"data"`
}

func TestMain(m *testing.M) {
	testConfig.ConfigureDb()

	err := testHelpers.SetUpTestDb()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	exitVal := m.Run()

	testHelpers.TearDownTestDb()

	os.Exit(exitVal)
}

func TestGetSalt(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"username":"%s"}`, testHelpers.ExampleGuest["email"]),
	}

	resp, err := getSaltHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("getSaltHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	saltData, err := deserializeBody(resp.Body)
	if saltData.Salt != testHelpers.ExampleCreds["salt"] || len(saltData.PrevSalts) != 0 || err != nil {
		t.Fatalf("getSaltHandler result %s/%d/%v, want %s/0/nil", saltData.Salt, len(saltData.PrevSalts), err, testHelpers.ExampleCreds["salt"])
	}
}

func TestGetSaltLocked(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"username":"%s"}`, testHelpers.ExampleGuest["email"]),
	}

	_ = testHelpers.LockAccount(testHelpers.ExampleGuest["email"])

	resp, err := getSaltHandler(context.TODO(), event)
	if resp.StatusCode != 429 || err != nil {
		t.Fatalf("getSaltHandler result %d/%v, want 429/nil", resp.StatusCode, err)
	}
}

func TestGetSaltBadData(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"username":"%s"}`, ""),
	}

	resp, err := getSaltHandler(context.TODO(), event)
	if resp.StatusCode != 400 || err != nil {
		t.Fatalf("getSaltHandler result %d/%v, want 400/nil", resp.StatusCode, err)
	}
}

func TestGetSaltNoUser(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"username":"%s"}`, "fake@test.fail"),
	}

	resp, err := getSaltHandler(context.TODO(), event)
	if resp.StatusCode != 500 || err != nil {
		t.Fatalf("getSaltHandler result %d/%v, want 500/nil", resp.StatusCode, err)
	}
}

func deserializeBody(body string) (SaltData, error) {
	var parsed DataBody

	b := []byte(body)
	err := json.Unmarshal(b, &parsed)

	return parsed.Data, err
}
