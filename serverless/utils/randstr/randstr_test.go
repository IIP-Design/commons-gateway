package randstr

import (
	"regexp"
	"testing"
)

const (
	STRLEN = 10
)

func TestRandDigitBytes(t *testing.T) {
	str, err := RandDigitBytes(STRLEN)
	want := regexp.MustCompile(`^\d+$`)
	if len(str) != STRLEN || !want.MatchString(str) || err != nil {
		t.Fatalf(`RandDigitBytes = %q, %v, want match for %#q, nil`, str, err, want)
	}
}

func TestRandStringBytes(t *testing.T) {
	str, err := RandStringBytes(STRLEN)
	want := regexp.MustCompile(`^[A-Za-z0-9_]+$`)
	if len(str) != STRLEN || !want.MatchString(str) || err != nil {
		t.Fatalf(`RandStringBytes = %q, %v, want match for %#q, nil`, str, err, want)
	}
}
