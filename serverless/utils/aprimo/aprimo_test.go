package aprimo

import (
	"bytes"
	"os"
	"testing"

	testConfig "github.com/IIP-Design/commons-gateway/test/config"
)

const (
	FILE_NAME        = "test-file.png"
	FILE_TYPE        = "image/png"
	FILE_DESCRIPTION = "Test File"
	APRIMO_TEAM      = "GPAVideo"
)

func TestAprimoSmallFileLifecycle(t *testing.T) {
	testConfig.ConfigureAprimo()

	token, err := GetAuthToken()
	if err != nil {
		t.Fatalf("GetAuthToken error %v", err)
	}

	fileData, err := os.ReadFile(FILE_NAME)
	if err != nil {
		t.Fatalf("ReadFile error %v", err)
	}

	fileToken, err := UploadFile(FILE_NAME, FILE_TYPE, bytes.NewBuffer(fileData), token)
	if err != nil {
		t.Fatalf("UploadFile error %v", err)
	}

	event := FileRecordInitEvent{
		Key:       FILE_NAME,
		FileType:  FILE_TYPE,
		FileToken: fileToken,
	}

	recordId, err := SubmitRecord(FILE_DESCRIPTION, event, APRIMO_TEAM, token)
	if err != nil {
		t.Fatalf("SubmitRecord error %v", err)
	}

	err = DeleteRecord(recordId, token)
	if err != nil {
		t.Fatalf("DeleteRecord error %v", err)
	}
}

func TestAprimoLargeFileLifecycle(t *testing.T) {
	testConfig.ConfigureAprimo()

	token, err := GetAuthToken()
	if err != nil {
		t.Fatalf("GetAuthToken error %v", err)
	}

	fileData, err := os.ReadFile(FILE_NAME)
	if err != nil {
		t.Fatalf("ReadFile error %v", err)
	}

	uri := InitFileUpload(FILE_NAME, token)
	success, err := UploadSegment(FILE_NAME, uri, &FileSegment{
		Segment:  0,
		FileType: FILE_TYPE,
		Data:     bytes.NewBuffer(fileData),
	}, token)

	if !success || err != nil {
		t.Fatalf("UploadSegment error %t, %v; want match fortrue, nil", success, err)
	}

	fileToken, err := CommitFileUpload(FILE_NAME, 1, uri, token)
	if err != nil {
		t.Fatalf("CommitFileUpload error %v", err)
	}

	event := FileRecordInitEvent{
		Key:       FILE_NAME,
		FileType:  FILE_TYPE,
		FileToken: fileToken,
	}

	recordId, err := SubmitRecord(FILE_DESCRIPTION, event, APRIMO_TEAM, token)
	if err != nil {
		t.Fatalf("SubmitRecord error %v", err)
	}

	err = DeleteRecord(recordId, token)
	if err != nil {
		t.Fatalf("DeleteRecord error %v", err)
	}
}
