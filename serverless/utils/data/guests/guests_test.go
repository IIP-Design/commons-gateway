package guests

import (
	"testing"
	"time"

	testHelpers "github.com/IIP-Design/commons-gateway/test/helpers"
)

func TestPwResetEveryOther(t *testing.T) {
	now := time.Now().Format(time.RFC3339)
	reset, err := shouldResetPassword(now, testHelpers.FarFutureDate(), false)

	if !reset || err != nil {
		t.Fatalf(`shouldResetPassword returned %t/%v, want true, nil`, reset, err)
	}
}

func TestPwResetAfterTime(t *testing.T) {
	now := time.Now().Format(time.RFC3339)
	reset, err := shouldResetPassword(now, testHelpers.FarFutureDate(), true)

	if !reset || err != nil {
		t.Fatalf(`shouldResetPassword returned %t/%v, want true, nil`, reset, err)
	}
}
