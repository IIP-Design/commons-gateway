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

func TestAdminDeactivate(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"username": testHelpers.ExampleAdmin["email"],
		},
	}

	resp, err := deactivateAdminHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("deactivateAdminHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	active, err := checkAdminDeactivated(testHelpers.ExampleAdmin["email"])
	if active || err != nil {
		t.Fatalf("checkAdminDeactivated result %t/%v, want false/nil", active, err)
	}
}

func checkAdminDeactivated(email string) (bool, error) {
	pool := data.ConnectToDB()
	defer pool.Close()

	query := `SELECT active FROM admins WHERE email = $1`
	row := pool.QueryRow(query, email)

	var active bool
	err := row.Scan(&active)

	return active, err
}
