package jwt

import (
	"regexp"
	"testing"
)

const (
	username   = "test@example.com"
	scope      = "guest"
	firstLogin = false
)

func TestGenerateJwt(t *testing.T) {
	token, err := GenerateJWT(username, scope, firstLogin)
	want := regexp.MustCompile(`^[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+$`)
	if !want.MatchString(token) || err != nil {
		t.Fatalf(`GenerateJWT = %q, %v, want match for %#q, nil`, token, err, want)
	}
}

func TestExtractBearerToken(t *testing.T) {
	token, _ := GenerateJWT(username, scope, firstLogin)
	bearer := "Bearer " + token
	token, err := extractBearerToken(bearer)
	want := regexp.MustCompile(`^[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+$`)
	if !want.MatchString(token) || err != nil {
		t.Fatalf(`extractBearerToken = %q, %v, want match for %#q, nil`, token, err, want)
	}
}

func TestExtractBearerTokenNoToken(t *testing.T) {
	bearer := ""
	token, err := extractBearerToken(bearer)
	if token != "" || err == nil {
		t.Fatalf(`extractBearerToken = %s, want error`, token)
	}
}

func TestVerifyToken(t *testing.T) {
	token, _ := GenerateJWT(username, scope, firstLogin)
	err := VerifyJWT(token, []string{scope})
	if err != nil {
		t.Fatalf(`VerifyJWT error %v, want nil`, err)
	}
}

func TestVerifyTokenBadScope(t *testing.T) {
	token, _ := GenerateJWT(username, scope, firstLogin)
	err := VerifyJWT(token, []string{"fail"})
	if err == nil {
		t.Fatal("VerifyJWT failed to generate an error")
	}
}

func TestCheckToken(t *testing.T) {
	token, _ := GenerateJWT(username, scope, firstLogin)
	bearer := "Bearer " + token
	err := CheckAuthToken(bearer, []string{scope})
	if err != nil {
		t.Fatalf(`CheckAuthToken error %v, want nil`, err)
	}
}

func TestExtractClientRole(t *testing.T) {
	token, _ := GenerateJWT(username, scope, firstLogin)
	bearer := "Bearer " + token
	role, err := ExtractClientRole(bearer)
	if role != scope || err != nil {
		t.Fatalf(`ExtractClientRole error %v, want nil`, err)
	}
}

func TestFormatJwt(t *testing.T) {
	fmt, err := FormatJWT(username, scope, firstLogin)
	want := regexp.MustCompile(`{\"token\":\"[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+\"}`)
	if !want.MatchString(fmt) || err != nil {
		t.Fatalf(`FormatJWT = %q, %v, want match for %#q, nil`, fmt, err, want)
	}
}
