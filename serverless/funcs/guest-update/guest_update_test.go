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
	FIRST_NAME = "Chandler"
	LAST_NAME  = "Cohen"
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

func TestUpdateGuestReal(t *testing.T) {
	event := makeGuestEvent(testHelpers.ExampleGuest["email"], testHelpers.ExampleTeam["id"])

	resp, err := guestUpdateHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("guestUpdateHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	firstName, lastName, err := checkGuestUpdated(testHelpers.ExampleGuest["email"])

	if firstName != FIRST_NAME || lastName != LAST_NAME || err != nil {
		t.Fatalf("Data is ill-formed: %s/%s/%v, want %s/%s/nil", firstName, lastName, err, FIRST_NAME, LAST_NAME)
	}
}

func TestUpdateGuestFakeUser(t *testing.T) {
	event := makeGuestEvent("fake@test.fail", testHelpers.ExampleTeam["id"])

	resp, err := guestUpdateHandler(context.TODO(), event)
	if resp.StatusCode != 404 || err != nil {
		t.Fatalf("guestUpdateHandler result %d/%v, want 404/nil", resp.StatusCode, err)
	}
}

func TestUpdateGuestFakeTeam(t *testing.T) {
	event := makeGuestEvent(testHelpers.ExampleGuest["email"], "ERROR")

	resp, err := guestUpdateHandler(context.TODO(), event)
	if resp.StatusCode != 404 || err != nil {
		t.Fatalf("guestUpdateHandler result %d/%v, want 404/nil", resp.StatusCode, err)
	}
}

func makeGuestEvent(email string, team string) events.APIGatewayProxyRequest {
	return events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"email":"%s","givenName":"%s","familyName":"%s","role":"%s","team":"%s"}`,
			email, FIRST_NAME, LAST_NAME, testHelpers.ExampleGuest["role"], team),
	}
}

func checkGuestUpdated(email string) (string, string, error) {
	pool := data.ConnectToDB()
	defer pool.Close()

	var firstName string
	var lastName string
	query := `SELECT first_name, last_name FROM guests WHERE email = $1`
	err := pool.QueryRow(query, email).Scan(&firstName, &lastName)

	return firstName, lastName, err
}
