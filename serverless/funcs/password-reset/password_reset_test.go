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

func TestResetPwReal(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"id": testHelpers.ExampleGuest["email"],
		},
	}

	resp, err := passwordResetHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("passwordResetHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	salt, hash, err := checkCredsUpdated(testHelpers.ExampleGuest["email"])

	if salt == testHelpers.ExampleCreds["salt"] || hash == testHelpers.ExampleCreds["pass_hash"] || err != nil {
		t.Fatalf("Data was not updated: %s/%s/%v, want %s/%s/nil", salt, hash, err, testHelpers.ExampleCreds["salt"], testHelpers.ExampleCreds["pass_hash"])
	}
}

func TestResetPwBadData(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"id": "",
		},
	}

	resp, err := passwordResetHandler(context.TODO(), event)
	if resp.StatusCode != 400 || err != nil {
		t.Fatalf("passwordResetHandler result %d/%v, want 400/nil", resp.StatusCode, err)
	}
}

func TestResetPwFakeUser(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"id": "fake@test.fail",
		},
	}

	resp, err := passwordResetHandler(context.TODO(), event)
	if resp.StatusCode != 404 || err != nil {
		t.Fatalf("passwordResetHandler result %d/%v, want 404/nil", resp.StatusCode, err)
	}
}

func checkCredsUpdated(email string) (string, string, error) {
	pool := data.ConnectToDB()
	defer pool.Close()

	var salt string
	var hash string
	query := `SELECT salt, pass_hash FROM invites WHERE invitee = $1`
	err := pool.QueryRow(query, email).Scan(&salt, &hash)

	return salt, hash, err
}
