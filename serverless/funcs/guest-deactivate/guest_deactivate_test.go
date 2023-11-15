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

	err := testHelpers.SetUpTestDb()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	exitVal := m.Run()

	testHelpers.TearDownTestDb()

	os.Exit(exitVal)
}

func TestDeactivateGuest(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"id": testHelpers.ExampleGuest["email"],
		},
	}

	resp, err := guestDeactivateHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("guestDeactivateHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	inactive, err := checkGuestDeactivated(testHelpers.ExampleGuest["email"])
	if !inactive || err != nil {
		t.Fatalf("checkGuestDeactivated result %t/%v, want false/nil", inactive, err)
	}
}

func TestMissDeactivation(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"id": "wrong@test.fail",
		},
	}

	resp, err := guestDeactivateHandler(context.TODO(), event)
	if resp.StatusCode != 404 || err != nil {
		t.Fatalf("guestDeactivateHandler result %d/%v, want 404/nil", resp.StatusCode, err)
	}
}

func checkGuestDeactivated(email string) (bool, error) {
	pool := data.ConnectToDB()
	defer pool.Close()

	var inactive bool
	query := `SELECT expiration < NOW() AS inactive FROM invites WHERE invitee = $1`
	err := pool.QueryRow(query, email).Scan(&inactive)

	return inactive, err
}
