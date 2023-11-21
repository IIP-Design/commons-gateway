package users

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

func TestGuest(t *testing.T) {
	_, success, err := CheckForExistingUser(testHelpers.ExampleGuest["email"], "guests")
	if !success || err != nil {
		t.Fatalf(`CheckForExistingUser result %t/%v, want true/nil`, success, err)
	}
}

func TestAdmin(t *testing.T) {
	_, success, err := CheckForExistingUser(testHelpers.ExampleAdmin["email"], "admins")
	if !success || err != nil {
		t.Fatalf(`CheckForExistingUser result %t/%v, want true/nil`, success, err)
	}
}

func TestMiss(t *testing.T) {
	_, success, err := CheckForExistingUser("fake@test.fail", "guests")
	if success || err != nil {
		t.Fatalf(`CheckForExistingUser result %t/%v, want true/nil`, success, err)
	}
}
