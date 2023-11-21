package teams

import (
	"fmt"
	"os"
	"testing"

	testConfig "github.com/IIP-Design/commons-gateway/test/config"
	testHelpers "github.com/IIP-Design/commons-gateway/test/helpers"
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

func TestGetByName(t *testing.T) {
	teamId, err := GetTeamIdByName(testHelpers.ExampleTeam["team_name"])
	if teamId != testHelpers.ExampleTeam["id"] || err != nil {
		t.Fatalf(`GetTeamIdByName returned %s/%v, want %s/nil`, teamId, err, testHelpers.ExampleTeam["id"])
	}
}
