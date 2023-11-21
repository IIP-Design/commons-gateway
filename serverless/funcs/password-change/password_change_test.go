package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	testConfig "github.com/IIP-Design/commons-gateway/test/config"
	testHelpers "github.com/IIP-Design/commons-gateway/test/helpers"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/randstr"
	"github.com/IIP-Design/commons-gateway/utils/security/hashing"
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/xid"
)

const (
	GOOD_PASSWORD = "goodPassword"
)

var prevPasswords = []string{"a", "b", "c", "d", "e"}

func TestMain(m *testing.M) {
	testConfig.ConfigureDb()

	err := testHelpers.SetUpTestDb()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	addPrevPasswords()

	exitVal := m.Run()

	cleanPasswords()
	testHelpers.TearDownTestDb()

	os.Exit(exitVal)
}

func TestBadUser(t *testing.T) {
	body, err := makeSubmission("fake@test.fail", GOOD_PASSWORD, testHelpers.ExampleCreds["pass_hash"])
	if err != nil {
		t.Fatalf(`makeSubmission error: %v`, err)
	}

	event := events.APIGatewayProxyRequest{
		Body: body,
	}

	resp, err := passwordChangeHandler(context.TODO(), event)
	if resp.StatusCode != 500 || err != nil {
		t.Fatalf("passwordChangeHandler result %d/%v, want 500/nil", resp.StatusCode, err)
	}
}

func TestBadPassword(t *testing.T) {
	body, err := makeSubmission(testHelpers.ExampleGuest["email"], GOOD_PASSWORD, "fail")
	if err != nil {
		t.Fatalf(`makeSubmission error: %v`, err)
	}

	event := events.APIGatewayProxyRequest{
		Body: body,
	}

	resp, err := passwordChangeHandler(context.TODO(), event)
	if resp.StatusCode != 500 || err != nil {
		t.Fatalf("passwordChangeHandler result %d/%v, want 500/nil", resp.StatusCode, err)
	}
}

func TestReusedPassword(t *testing.T) {
	body, err := makeSubmission(testHelpers.ExampleGuest["email"], prevPasswords[0], testHelpers.ExampleCreds["pass_hash"])
	if err != nil {
		t.Fatalf(`makeSubmission error: %v`, err)
	}

	event := events.APIGatewayProxyRequest{
		Body: body,
	}

	resp, err := passwordChangeHandler(context.TODO(), event)
	if resp.StatusCode != 409 || err != nil {
		t.Fatalf("passwordChangeHandler result %d/%v, want 409/nil", resp.StatusCode, err)
	}
}

func TestUpdateSuccess(t *testing.T) {
	body, err := makeSubmission(testHelpers.ExampleGuest["email"], GOOD_PASSWORD, testHelpers.ExampleCreds["pass_hash"])
	if err != nil {
		t.Fatalf(`makeSubmission error: %v`, err)
	}

	event := events.APIGatewayProxyRequest{
		Body: body,
	}

	resp, err := passwordChangeHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("passwordChangeHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}
}

func addPrevPasswords() error {
	pool := data.ConnectToDB()
	defer pool.Close()

	query := "INSERT INTO password_history ( id, user_id, creation_date, salt, pass_hash ) VALUES ( $1, $2, NOW(), $3, $4)"

	for _, pass := range prevPasswords {
		id := xid.New()
		salt, _ := randstr.RandStringBytes(10)
		hash := hashing.GenerateHash(pass, salt)

		_, err := pool.Exec(query, id, testHelpers.ExampleGuest["email"], salt, hash)
		if err != nil {
			return err
		}
	}

	return nil
}

func getPrevSaltHashes(password string) ([]string, error) {
	var hashes []string
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	query := "SELECT salt FROM password_history WHERE user_id = $1"
	rows, err := pool.Query(query, testHelpers.ExampleGuest["email"])

	if err != nil {
		return hashes, err
	}

	defer rows.Close()

	for rows.Next() {
		var salt string

		if err = rows.Scan(&salt); err != nil {
			return hashes, err
		}

		hashes = append(hashes, hashing.GenerateHash(password, salt))
	}

	return hashes, err
}

func makeSubmission(email string, newPassword string, currPassword string) (string, error) {
	salt, _ := randstr.RandStringBytes(10)
	hash := hashing.GenerateHash(newPassword, salt)

	prevHashes, err := getPrevSaltHashes(newPassword)
	if err != nil {
		return "", err
	}

	sub := PasswordReset{
		CurrentPasswordHash: currPassword,
		Email:               email,
		HashedPriorSalts:    prevHashes,
		NewPasswordHash:     hash,
		NewSalt:             salt,
	}

	ret, err := json.Marshal(sub)
	return string(ret), err
}

func cleanPasswords() error {
	pool := data.ConnectToDB()
	defer pool.Close()

	query := "DELETE FROM password_history WHERE user_id = $1"
	_, err := pool.Exec(query, testHelpers.ExampleGuest["email"])
	return err
}
