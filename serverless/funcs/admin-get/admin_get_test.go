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

	err := testHelpers.SetUpTestDb()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	exitVal := m.Run()

	testHelpers.TearDownTestDb()

	os.Exit(exitVal)
}

func TestGetAdmin(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"username": testHelpers.ExampleAdmin["email"],
		},
	}

	resp, err := getAdminHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("getAdminHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}
}

func TestMissAdmin(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"username": "wrong@test.fail",
		},
	}

	resp, err := getAdminHandler(context.TODO(), event)
	if resp.StatusCode != 404 || err != nil {
		t.Fatalf("getAdminHandler result %d/%v, want 404/nil", resp.StatusCode, err)
	}
}
