package main

// TODO: Ensure test success is not order-dependent

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	testConfig "github.com/IIP-Design/commons-gateway/test/config"
	testHelpers "github.com/IIP-Design/commons-gateway/test/helpers"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/aws/aws-lambda-go/events"
)

const (
	REQUEST_ID = "9m4e2mr0ui3e8request"
	CODE       = "123456"
)

func TestMain(m *testing.M) {
	testConfig.ConfigureDb()

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

	cleanUpMfa()
	testHelpers.CleanupInvites(testHelpers.ExampleGuest2["email"])
	testHelpers.TearDownTestDb()

	os.Exit(exitVal)
}

func TestBadPassword(t *testing.T) {
	addMfa()

	event := events.APIGatewayProxyRequest{
		Body: makeJsonBody("fail", testHelpers.ExampleGuest["email"], CODE),
	}

	attempts1, err := checkLoginAttempts(testHelpers.ExampleGuest["email"])
	if err != nil {
		t.Fatalf("checkLoginAttempts error %v", err)
	}

	resp, err := authenticationHandler(context.TODO(), event)
	if resp.StatusCode != 403 || err != nil {
		t.Fatalf("authenticationHandler result %d/%v, want 403/nil", resp.StatusCode, err)
	}

	attempts2, err := checkLoginAttempts(testHelpers.ExampleGuest["email"])
	if err != nil {
		t.Fatalf("checkLoginAttempts error %v", err)
	}

	if attempts2 <= attempts1 {
		t.Fatal("Did not record failed login attempt")
	}
}

func TestBad2fa(t *testing.T) {
	addMfa()

	event := events.APIGatewayProxyRequest{
		Body: makeJsonBody(testHelpers.ExampleCreds["pass_hash"], testHelpers.ExampleGuest["email"], "fail"),
	}

	attempts1, err := checkLoginAttempts(testHelpers.ExampleGuest["email"])
	if err != nil {
		t.Fatalf("checkLoginAttempts error %v", err)
	}

	resp, err := authenticationHandler(context.TODO(), event)
	if resp.StatusCode != 403 || err != nil {
		t.Fatalf("authenticationHandler result %d/%v, want 403/nil", resp.StatusCode, err)
	}

	attempts2, err := checkLoginAttempts(testHelpers.ExampleGuest["email"])
	if err != nil {
		t.Fatalf("checkLoginAttempts error %v", err)
	}

	if attempts2 <= attempts1 {
		t.Fatal("Did not record failed login attempt")
	}
}

func TestPending(t *testing.T) {
	addMfa()

	testHelpers.AddPendingGuest()
	event := events.APIGatewayProxyRequest{
		Body: makeJsonBody(testHelpers.ExampleCreds["pass_hash"], testHelpers.ExampleGuest2["email"], CODE),
	}

	attempts1, err := checkLoginAttempts(testHelpers.ExampleGuest2["email"])
	if err != nil {
		t.Fatalf("checkLoginAttempts error %v", err)
	}

	resp, err := authenticationHandler(context.TODO(), event)
	if resp.StatusCode != 403 || err != nil {
		t.Fatalf("authenticationHandler result %d/%v, want 403/nil", resp.StatusCode, err)
	}

	attempts2, err := checkLoginAttempts(testHelpers.ExampleGuest2["email"])
	if err != nil {
		t.Fatalf("checkLoginAttempts error %v", err)
	}

	if attempts2 <= attempts1 {
		t.Fatal("Did not record failed login attempt")
	}
}

func TestSuccess(t *testing.T) {
	addMfa()

	event := events.APIGatewayProxyRequest{
		Body: makeJsonBody(testHelpers.ExampleCreds["pass_hash"], testHelpers.ExampleGuest["email"], CODE),
	}

	resp, err := authenticationHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("authenticationHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	attempts, err := checkLoginAttempts(testHelpers.ExampleGuest["email"])
	if attempts > 0 || err != nil {
		t.Fatalf("checkLoginAttempts error or not reset: %d/%v", attempts, err)
	}
}

func TestLocked(t *testing.T) {
	addMfa()

	testHelpers.LockAccount(testHelpers.ExampleGuest["email"])
	event := events.APIGatewayProxyRequest{
		Body: makeJsonBody(testHelpers.ExampleCreds["pass_hash"], testHelpers.ExampleGuest["email"], CODE),
	}

	resp, err := authenticationHandler(context.TODO(), event)
	if resp.StatusCode != 429 || err != nil {
		t.Fatalf("authenticationHandler result %d/%v, want 429/nil", resp.StatusCode, err)
	}
}

func TestExpired(t *testing.T) {
	addMfa()
	testHelpers.DeactivateGuest(testHelpers.ExampleGuest["email"])

	event := events.APIGatewayProxyRequest{
		Body: makeJsonBody(testHelpers.ExampleCreds["pass_hash"], testHelpers.ExampleGuest["email"], CODE),
	}

	attempts1, err := checkLoginAttempts(testHelpers.ExampleGuest["email"])
	if err != nil {
		t.Fatalf("checkLoginAttempts error %v", err)
	}

	resp, err := authenticationHandler(context.TODO(), event)
	if resp.StatusCode != 403 || err != nil {
		t.Fatalf("authenticationHandler result %d/%v, want 403/nil", resp.StatusCode, err)
	}

	attempts2, err := checkLoginAttempts(testHelpers.ExampleGuest["email"])
	if err != nil {
		t.Fatalf("checkLoginAttempts error %v", err)
	}

	if attempts2 <= attempts1 {
		t.Fatal("Did not record failed login attempt")
	}
}

func makeJsonBody(hash string, email string, code string) string {
	return fmt.Sprintf(`{
		"hash": "%s",
		"username": "%s",
		"mfa": {
			"code": "%s",
			"id": "%s"
		}
	}`,
		hash, email, code, REQUEST_ID)
}

func checkLoginAttempts(email string) (int, error) {
	pool := data.ConnectToDB()
	defer pool.Close()

	var attempts int
	query := `SELECT login_attempt FROM guests WHERE email = $1`
	err := pool.QueryRow(query, email).Scan(&attempts)

	return attempts, err
}

func addMfa() error {
	pool := data.ConnectToDB()
	defer pool.Close()

	currentTime := time.Now()

	insertMfa := `INSERT INTO mfa( request_id, code, date_created ) VALUES ( $1, $2, $3 ) ON CONFLICT ON CONSTRAINT mfa_pkey DO NOTHING;`
	_, err := pool.Exec(insertMfa, REQUEST_ID, CODE, currentTime)

	return err
}

func cleanUpMfa() error {
	pool := data.ConnectToDB()
	defer pool.Close()

	query := `DELETE FROM mfa WHERE code = $1`
	_, err := pool.Exec(query, CODE)

	return err
}
