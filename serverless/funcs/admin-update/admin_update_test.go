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

const (
	GIVEN_NAME  = "Dawn"
	FAMILY_NAME = "Schafer"
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

func TestUpdateAdmin(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"email":"%s","givenName":"%s","familyName":"%s","role":"%s","team":"%s", "active":false}`,
			testHelpers.ExampleAdmin["email"], GIVEN_NAME, FAMILY_NAME, testHelpers.ExampleAdmin["role"], testHelpers.ExampleTeam["id"]),
	}

	resp, err := updateAdminHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("updateAdminHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	pool := data.ConnectToDB()
	defer pool.Close()

	query := "SELECT first_name, last_name FROM admins WHERE email = $1"

	var firstName string
	var lastName string

	pool.QueryRow(query, testHelpers.ExampleAdmin["email"]).Scan(&firstName, &lastName)

	if firstName != GIVEN_NAME || lastName != FAMILY_NAME {
		t.Fatalf("Data is %s/%s, want %s/%s", firstName, lastName, GIVEN_NAME, FAMILY_NAME)
	}
}

func TestUpdateFakeAdmin(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"email":"%s","givenName":"%s","familyName":"%s","role":"%s","team":"%s", "active":false}`,
			"wrong@test.fail", GIVEN_NAME, FAMILY_NAME, testHelpers.ExampleAdmin["role"], testHelpers.ExampleTeam["id"]),
	}

	resp, err := updateAdminHandler(context.TODO(), event)
	if resp.StatusCode != 404 || err != nil {
		t.Fatalf("updateAdminHandler result %d/%v, want 404/nil", resp.StatusCode, err)
	}
}

func TestUpdateFakeTeam(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"email":"%s","givenName":"%s","familyName":"%s","role":"%s","team":"%s", "active":false}`,
			testHelpers.ExampleAdmin["email"], GIVEN_NAME, FAMILY_NAME, testHelpers.ExampleAdmin["role"], "ERROR"),
	}

	resp, err := updateAdminHandler(context.TODO(), event)
	if resp.StatusCode != 404 || err != nil {
		t.Fatalf("updateAdminHandler result %d/%v, want 404/nil", resp.StatusCode, err)
	}
}
