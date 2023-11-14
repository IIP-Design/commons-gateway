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
	GIVEN_NAME  = "Carmen"
	FAMILY_NAME = "Lowell"
	EMAIL       = "new.admin@example.com"
	ROLE        = "admin"
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

	cleanupadmins()
	testHelpers.TearDownTestDb()

	os.Exit(exitVal)
}

func TestCreateAdmin(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"email":"%s","givenName":"%s","familyName":"%s","role":"%s","team":"%s"}`,
			EMAIL, GIVEN_NAME, FAMILY_NAME, ROLE, testHelpers.ExampleTeam["id"]),
	}

	resp, err := newAdminHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("newAdminHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	pool := data.ConnectToDB()
	defer pool.Close()

	query := "SELECT first_name, last_name, role FROM admins WHERE email = $1"

	var firstName string
	var lastName string
	var role string

	pool.QueryRow(query, EMAIL).Scan(&firstName, &lastName, &role)

	if role != ROLE || firstName != GIVEN_NAME || lastName != FAMILY_NAME {
		t.Fatalf("Data is %s/%s/%s, want %s/%s/%s", role, firstName, lastName, ROLE, GIVEN_NAME, FAMILY_NAME)
	}
}

func TestCreateExistingAdmin(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"email":"%s","givenName":"%s","familyName":"%s","role":"%s","team":"%s"}`,
			EMAIL, GIVEN_NAME, FAMILY_NAME, ROLE, testHelpers.ExampleTeam["id"]),
	}

	resp, err := newAdminHandler(context.TODO(), event)
	if resp.StatusCode != 409 || err != nil {
		t.Fatalf("newAdminHandler result %d/%v, want 409/nil", resp.StatusCode, err)
	}
}

func cleanupadmins() {
	pool := data.ConnectToDB()
	defer pool.Close()

	query := "DELETE FROM admins WHERE email = $1"
	pool.Exec(query, EMAIL)
}
