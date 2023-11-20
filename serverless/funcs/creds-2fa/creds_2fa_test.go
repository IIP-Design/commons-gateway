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

const (
	QUEUE_NAME = "email_test_queue"
	ENV        = "EMAIL_MFA_QUEUE"
)

type MfaEntry struct {
	Id string `json:"requestId"`
}

type DataBody struct {
	Data MfaEntry `json:"data"`
}

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

	os.Exit(exitVal)
}

func TestRegMfaNoQueue(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"username": testHelpers.ExampleGuest["email"],
		},
	}

	resp, err := generateMfaHandler(context.TODO(), event)
	if resp.StatusCode != 500 || err != nil {
		t.Fatalf("generateMfaHandler result %d/%v, want 500/nil", resp.StatusCode, err)
	}
}

func TestRegMfa(t *testing.T) {
	_, client, err := testHelpers.CreateTestQueue(QUEUE_NAME, ENV)
	if err != nil {
		t.Fatalf("CreateTestQueue error %v", err)
	}

	event := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"username": testHelpers.ExampleGuest["email"],
		},
	}

	resp, err := generateMfaHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("generateMfaHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	testHelpers.DeleteQueue(QUEUE_NAME, client)

	mfa, err := deserializeBody(resp.Body)
	if mfa.Id == "" || err != nil {
		t.Fatalf("generateMfaHandler result empty or error %v", err)
	}

	code, err := checkMfaRegistered(mfa.Id)
	if code == "" || err != nil {
		t.Fatalf("Code result empty or error %v", err)
	}
}

func TestMissMfa(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"username": "wrong@test.fail",
		},
	}

	resp, err := generateMfaHandler(context.TODO(), event)
	if resp.StatusCode != 404 || err != nil {
		t.Fatalf("generateMfaHandler result %d/%v, want 404/nil", resp.StatusCode, err)
	}
}

func TestBadData(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{},
	}

	resp, err := generateMfaHandler(context.TODO(), event)
	if resp.StatusCode != 400 || err != nil {
		t.Fatalf("generateMfaHandler result %d/%v, want 400/nil", resp.StatusCode, err)
	}
}

func deserializeBody(body string) (MfaEntry, error) {
	var parsed DataBody

	b := []byte(body)
	err := json.Unmarshal(b, &parsed)

	return parsed.Data, err
}

func checkMfaRegistered(requestId string) (string, error) {
	pool := data.ConnectToDB()
	defer pool.Close()

	var code string
	query := `SELECT code FROM mfa WHERE request_id = $1`
	err := pool.QueryRow(query, requestId).Scan(&code)

	return code, err
}
