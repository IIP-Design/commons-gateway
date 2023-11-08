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

	teamId, err := testHelpers.SetUpTestDb()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	exitVal := m.Run()

	testHelpers.TearDownTestDb(teamId)

	os.Exit(exitVal)
}

func TestGetAdmin(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"username": "admin@example.com",
		},
	}

	resp, err := getAdminHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("getAdminHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}
}
