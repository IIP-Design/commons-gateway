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
	S3_ID       = "TESTID"
	FILE_TYPE   = "image/png"
	DESCRIPTION = "Image"
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

	cleanupMetadataRecords()
	testHelpers.TearDownTestDb()

	os.Exit(exitVal)
}

func TestUploadMetadata(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"key":"%s","email":"%s","team":"%s","fileType":"%s","description":"%s"}`,
			S3_ID, testHelpers.ExampleAdmin["email"], testHelpers.ExampleTeam["id"], FILE_TYPE, DESCRIPTION),
	}

	resp, err := newUploadHandler(context.TODO(), event)
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("newUploadHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}

	pool := data.ConnectToDB()
	defer pool.Close()

	query := "SELECT s3_id, team_id, file_type FROM uploads WHERE s3_id = $1"

	var s3Id string
	var team string
	var fileType string

	pool.QueryRow(query, S3_ID).Scan(&s3Id, &team, &fileType)

	if s3Id != S3_ID || team != testHelpers.ExampleTeam["id"] || fileType != FILE_TYPE {
		t.Fatalf("Data is %s/%s/%s, want %s/%s/%s", s3Id, team, fileType, S3_ID, testHelpers.ExampleTeam["id"], FILE_TYPE)
	}
}

func cleanupMetadataRecords() {
	pool := data.ConnectToDB()
	defer pool.Close()

	query := "DELETE FROM uploads WHERE s3_id = $1"
	pool.Exec(query, S3_ID)
}
