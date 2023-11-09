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

func TestGetAdmin(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"team":"%s"}`, testHelpers.ExampleTeam["id"]),
	}

	resp, err := getUploaderHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("getUploaderHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	result, err := testHelpers.DeserializeBodyArray(resp.Body)
	if err != nil {
		t.Fatalf("DeserializeBodyArray result %v, want nil", err)
	}

	if len(result) == 0 || result[0].(map[string]any)["email"] != testHelpers.ExampleGuest["email"] {
		t.Fatal("Body has no results or is ill-formed")
	}
}
