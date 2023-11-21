package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	testConfig "github.com/IIP-Design/commons-gateway/test/config"
	testHelpers "github.com/IIP-Design/commons-gateway/test/helpers"
	"github.com/aws/aws-lambda-go/events"
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

func TestGetAdmins(t *testing.T) {
	event := events.APIGatewayProxyRequest{}

	resp, err := getAdminsHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("getAdminsHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	body := resp.Body
	arrayRe := regexp.MustCompile(`(?m)^{"data":\[{`)   // An array was returned
	emailRe := regexp.MustCompile(`admin@example\.com`) // The test email is there

	if !arrayRe.MatchString(body) || !emailRe.MatchString(body) {
		t.Fatalf("Data is ill-formed: %s", body)
	}
}
