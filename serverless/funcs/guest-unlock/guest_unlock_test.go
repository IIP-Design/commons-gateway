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
	MESSAGE_ID = "98765"
)

func TestMain(m *testing.M) {
	testConfig.ConfigureDb()

	err := testHelpers.SetUpTestDb()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	testHelpers.AddPendingGuest()
	testHelpers.LockAccount(testHelpers.ExampleGuest["email"])
	testHelpers.LockAccount(testHelpers.ExampleGuest2["email"])

	exitVal := m.Run()

	testHelpers.TearDownTestDb()
	testHelpers.CleanupInvites(testHelpers.ExampleGuest2["email"])

	os.Exit(exitVal)
}

func TestUnlockMiss(t *testing.T) {
	eventBody := fmt.Sprintf(`{"username":"%s"}`, "fake@test.fail")
	event := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId: MESSAGE_ID,
				Body:      eventBody,
			},
		},
	}

	resp, err := unlockGuestHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("unlockGuestHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	g1Locked, err := checkGuestLocked(testHelpers.ExampleGuest["email"])
	if !g1Locked || err != nil {
		t.Fatalf("checkGuestLocked result %t/%v, want true/nil", g1Locked, err)
	}

	g2Locked, err := checkGuestLocked(testHelpers.ExampleGuest2["email"])
	if !g2Locked || err != nil {
		t.Fatalf("checkGuestLocked result %t/%v, want true/nil", g2Locked, err)
	}
}

func TestUnlock(t *testing.T) {
	eventBody := fmt.Sprintf(`{"username":"%s"}`, testHelpers.ExampleGuest["email"])
	event := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId: MESSAGE_ID,
				Body:      eventBody,
			},
		},
	}

	resp, err := unlockGuestHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("unlockGuestHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	g1Locked, err := checkGuestLocked(testHelpers.ExampleGuest["email"])
	if g1Locked || err != nil {
		t.Fatalf("checkGuestLocked result %t/%v, want false/nil", g1Locked, err)
	}

	g2Locked, err := checkGuestLocked(testHelpers.ExampleGuest2["email"])
	if !g2Locked || err != nil {
		t.Fatalf("checkGuestLocked result %t/%v, want true/nil", g2Locked, err)
	}
}

func checkGuestLocked(email string) (bool, error) {
	pool := data.ConnectToDB()
	defer pool.Close()

	var locked bool
	query := `SELECT locked FROM guests WHERE email = $1`
	err := pool.QueryRow(query, email).Scan(&locked)

	return locked, err
}
