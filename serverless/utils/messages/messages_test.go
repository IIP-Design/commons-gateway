package messages

import (
	"errors"
	"testing"
)

func TestSuccesMessage(t *testing.T) {
	resp, err := SendSuccessMessage()
	if err != nil {
		t.Fatalf(`SendSuccessMessage error %v, want nil`, err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf(`SendSuccessMessage status %d, want %d`, resp.StatusCode, 200)
	}

	if resp.Body != `{"message":"success"}` {
		t.Fatalf(`SendSuccessMessage body %s, want %s`, resp.Body, `{"message":"success"}`)
	}
}

func TestServerError(t *testing.T) {
	testErr := errors.New("TEST")

	resp, err := SendServerError(testErr)
	if err != nil {
		t.Fatalf(`SendServerError error %v, want nil`, err)
	}

	if resp.StatusCode != 500 {
		t.Fatalf(`SendServerError status %d, want %d`, resp.StatusCode, 500)
	}

	if resp.Body != testErr.Error() {
		t.Fatalf(`SendCustomError body %s, want %s`, resp.Body, testErr.Error())
	}
}

func TestCustomError(t *testing.T) {
	resp, err := SendCustomError(nil, 409)
	if err != nil {
		t.Fatalf(`SendCustomError error %v, want nil`, err)
	}

	if resp.StatusCode != 409 {
		t.Fatalf(`SendCustomError status %d, want %d`, resp.StatusCode, 409)
	}

	if resp.Body != `{"error":"resource conflict"}` {
		t.Fatalf(`SendCustomError body %s, want %s`, resp.Body, `{"error":"resource conflict"}`)
	}
}
