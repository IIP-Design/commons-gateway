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
	TEAM_NAME = "Aftermath"
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

func TestUpdateTeamName(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"team":"%s","teamName":"%s", "teamAprimo":"%s", "active":%t}`,
			testHelpers.ExampleTeam["id"], TEAM_NAME, testHelpers.ExampleTeam["aprimo_name"], true),
	}

	resp, err := teamUpdateHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("teamUpdateHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	teamName, active, err := checkTeamUpdated()
	if teamName != TEAM_NAME || !active || err != nil {
		t.Fatalf("Team was not updated: %s/%t/%v", teamName, active, err)
	}
}

func TestUpdateTeamActive(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"team":"%s","active":%t}`,
			testHelpers.ExampleTeam["id"], false),
	}

	resp, err := teamUpdateHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("teamUpdateHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	_, active, err := checkTeamUpdated()
	if active || err != nil {
		t.Fatalf("Team was not deactivated: %t/%v", active, err)
	}
}

func TestUpdateTeamBadData(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"team":"%s","teamName":"%s", "teamAprimo":"%s", "active":%t}`,
			"", TEAM_NAME, testHelpers.ExampleTeam["aprimo_name"], true),
	}

	resp, err := teamUpdateHandler(context.TODO(), event)
	if resp.StatusCode != 400 || err != nil {
		t.Fatalf("teamUpdateHandler result %d/%v, want 400/nil", resp.StatusCode, err)
	}
}

func TestUpdateTeamMiss(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"team":"%s","teamName":"%s", "teamAprimo":"%s", "active":%t}`,
			"ERROR", TEAM_NAME, testHelpers.ExampleTeam["aprimo_name"], true),
	}

	resp, err := teamUpdateHandler(context.TODO(), event)
	if resp.StatusCode != 404 || err != nil {
		t.Fatalf("teamUpdateHandler result %d/%v, want 404/nil", resp.StatusCode, err)
	}
}

func checkTeamUpdated() (string, bool, error) {
	pool := data.ConnectToDB()
	defer pool.Close()

	var teamName string
	var active bool
	query := `SELECT team_name, active FROM teams WHERE id = $1`
	err := pool.QueryRow(query, testHelpers.ExampleTeam["id"]).Scan(&teamName, &active)

	return teamName, active, err
}
