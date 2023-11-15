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
	"github.com/IIP-Design/commons-gateway/utils/types"
	"github.com/aws/aws-lambda-go/events"
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

	testHelpers.TearDownTestDb()

	os.Exit(exitVal)
}

func TestGetTeams(t *testing.T) {
	event := events.APIGatewayProxyRequest{}

	resp, err := getTeamsHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("getTeamsHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	body := resp.Body
	teams, err := deserializeBody(body)
	idx := slices.IndexFunc(teams, func(t types.Team) bool { return t.Id == testHelpers.ExampleTeam["id"] })

	if err != nil || len(teams) == 0 || idx == -1 {
		t.Fatalf("Data is ill-formed: %v/%d/%d", err, len(teams), idx)
	}
}

func deserializeBody(body string) ([]types.Team, error) {
	var parsed DataBody

	b := []byte(body)
	err := json.Unmarshal(b, &parsed)

	return parsed.Data, err
}
