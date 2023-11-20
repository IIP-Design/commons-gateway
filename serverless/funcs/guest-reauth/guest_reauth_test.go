package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	testConfig "github.com/IIP-Design/commons-gateway/test/config"
	testHelpers "github.com/IIP-Design/commons-gateway/test/helpers"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/security/jwt"
	"github.com/aws/aws-lambda-go/events"
)

func TestMain(m *testing.M) {
	testConfig.ConfigureDb()
	testConfig.ConfigureEmail()

	testHelpers.TearDownTestDb()
	err := testHelpers.SetUpTestDb()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = testHelpers.AddPendingGuest()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	exitVal := m.Run()

	testHelpers.CleanupInvites(testHelpers.ExampleGuest2["email"])
	testHelpers.TearDownTestDb()

	os.Exit(exitVal)
}

func TestBadScope(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body:    makeJsonBody("fake@test.fail", testHelpers.ExampleAdmin["email"]),
		Headers: map[string]string{},
	}

	resp, err := guestReauthHandler(context.TODO(), event)
	if resp.StatusCode != 500 || err != nil {
		t.Fatalf("guestReauthHandler result %d/%v, want 500/nil", resp.StatusCode, err)
	}
}

func TestBadUser(t *testing.T) {
	token, _ := jwt.GenerateJWT(testHelpers.ExampleAdmin["email"], "admin", false)

	event := events.APIGatewayProxyRequest{
		Body: makeJsonBody("fake@test.fail", testHelpers.ExampleAdmin["email"]),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf(`Bearer %s`, token),
		},
	}

	resp, err := guestReauthHandler(context.TODO(), event)
	if resp.StatusCode != 404 || err != nil {
		t.Fatalf("guestReauthHandler result %d/%v, want 404/nil", resp.StatusCode, err)
	}
}

func TestPendingUser(t *testing.T) {
	token, _ := jwt.GenerateJWT(testHelpers.ExampleAdmin["email"], "admin", false)

	event := events.APIGatewayProxyRequest{
		Body: makeJsonBody(testHelpers.ExampleGuest2["email"], testHelpers.ExampleAdmin["email"]),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf(`Bearer %s`, token),
		},
	}

	resp, err := guestReauthHandler(context.TODO(), event)
	if resp.StatusCode != 409 || err != nil {
		t.Fatalf("guestReauthHandler result %d/%v, want 409/nil", resp.StatusCode, err)
	}
}

func TestActiveUser(t *testing.T) {
	token, _ := jwt.GenerateJWT(testHelpers.ExampleAdmin["email"], "admin", false)

	event := events.APIGatewayProxyRequest{
		Body: makeJsonBody(testHelpers.ExampleGuest["email"], testHelpers.ExampleAdmin["email"]),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf(`Bearer %s`, token),
		},
	}

	resp, err := guestReauthHandler(context.TODO(), event)
	if resp.StatusCode != 409 || err != nil {
		t.Fatalf("guestReauthHandler result %d/%v, want 409/nil", resp.StatusCode, err)
	}
}

func TestUserGuestAdmin(t *testing.T) {
	approveGuest(testHelpers.ExampleGuest2["email"])
	testHelpers.DeactivateGuest(testHelpers.ExampleGuest2["email"])

	token, _ := jwt.GenerateJWT(testHelpers.ExampleGuest["email"], "guest admin", false)

	event := events.APIGatewayProxyRequest{
		Body: makeJsonBody(testHelpers.ExampleGuest2["email"], testHelpers.ExampleGuest["email"]),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf(`Bearer %s`, token),
		},
	}

	resp, err := guestReauthHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("guestReauthHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	active, err := checkUserReauth(testHelpers.ExampleGuest2["email"])
	if !active || err != nil {
		t.Fatalf("checkUserReauth result %t/%v, want true/nil", active, err)
	}
}

func TestUserAdmin(t *testing.T) {
	testHelpers.DeactivateGuest(testHelpers.ExampleGuest["email"])
	token, _ := jwt.GenerateJWT(testHelpers.ExampleAdmin["email"], "admin", false)

	event := events.APIGatewayProxyRequest{
		Body: makeJsonBody(testHelpers.ExampleGuest["email"], testHelpers.ExampleAdmin["email"]),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf(`Bearer %s`, token),
		},
	}

	resp, err := guestReauthHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("guestReauthHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	active, err := checkUserReauth(testHelpers.ExampleGuest["email"])
	if !active || err != nil {
		t.Fatalf("checkUserReauth result %t/%v, want true/nil", active, err)
	}
}

func makeJsonBody(email string, admin string) string {
	return fmt.Sprintf(`{
		"expiration": "2025-12-01T10:00:00Z",
		"email": "%s",
		"admin": "%s"
	}`, email, admin)
}

func approveGuest(email string) error {
	pool := data.ConnectToDB()
	defer pool.Close()

	query := `UPDATE invites SET pending = FALSE WHERE invitee = $1`
	_, err := pool.Exec(query, email)

	return err
}

func checkUserReauth(email string) (bool, error) {
	pool := data.ConnectToDB()
	defer pool.Close()

	var active bool
	query := `SELECT expiration >= NOW() AS active FROM invites WHERE invitee = $1 ORDER BY date_invited DESC LIMIT 1;`
	err := pool.QueryRow(query, email).Scan(&active)

	return active, err
}
