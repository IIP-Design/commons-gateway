package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"testing"

	testConfig "github.com/IIP-Design/commons-gateway/test/config"
	testHelpers "github.com/IIP-Design/commons-gateway/test/helpers"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/types"
	"github.com/aws/aws-lambda-go/events"
)

const (
	TEAM_NAME   = "Aftermath"
	TEAM_APRIMO = "AM.Video"
)

type DataBody struct {
	Data []types.Team `json:"data"`
}

func TestMain(m *testing.M) {
	testConfig.ConfigureDb()

	testHelpers.TearDownTestDb()
	err := testHelpers.SetUpTestDb()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	exitVal := m.Run()

	cleanupTeams()
	testHelpers.TearDownTestDb()

	os.Exit(exitVal)
}

func TestCreateTeam(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"teamName":"%s", "teamAprimo":"%s"}`,
			TEAM_NAME, TEAM_APRIMO),
	}

	resp, err := newTeamHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("newTeamHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	body := resp.Body
	teams, err := deserializeBody(body)
	idx := slices.IndexFunc(teams, func(t types.Team) bool { return t.Name == TEAM_NAME })

	if err != nil || idx == -1 {
		t.Fatalf("Data is ill-formed: %v/%d", err, idx)
	}
}

func TestCreateTeamConflict(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"teamName":"%s", "teamAprimo":"%s"}`,
			testHelpers.ExampleTeam["team_name"], TEAM_APRIMO),
	}

	resp, err := newTeamHandler(context.TODO(), event)
	if resp.StatusCode != 409 || err != nil {
		t.Fatalf("newTeamHandler result %d/%v, want 409/nil", resp.StatusCode, err)
	}
}

func TestCreateTeamBadData(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"teamName":"%s", "teamAprimo":"%s"}`,
			"", ""),
	}

	resp, err := newTeamHandler(context.TODO(), event)
	if resp.StatusCode != 400 || err != nil {
		t.Fatalf("newTeamHandler result %d/%v, want 400/nil", resp.StatusCode, err)
	}
}

func deserializeBody(body string) ([]types.Team, error) {
	var parsed DataBody

	b := []byte(body)
	err := json.Unmarshal(b, &parsed)

	return parsed.Data, err
}

func cleanupTeams() {
	pool := data.ConnectToDB()
	defer pool.Close()

	query := "DELETE FROM teams WHERE team_name = $1"
	pool.Exec(query, TEAM_NAME)
}
