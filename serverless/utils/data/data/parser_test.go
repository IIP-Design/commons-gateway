package data

import (
	"testing"
)

func makeJsonBody() string {
	return `{
		"email": "email@example.com",
		"expiration": "2023-12-01T12:00:00Z",
		"givenName": "John",
		"familyName": "Public",
		"role": "guest",
		"team": "abcde",
		"teamName": "GPA Video",
		"teamAprimo": "GPAVideo",
	
		"active": true,
		"hash": "hash",
		"invitee": {
			"email": "invitee@example.com",
			"expiration": "2023-12-01",
			"givenName": "Carmen",
			"familyName": "Lowell",
			"role": "guest",
			"team": "abcde",
			"teamName": "GPA Video",
			"teamAprimo": "GPAVideo"
		},
		"inviter": "inviter@example.com",
		"admin": "admin@example.com",
		"mfa": {
			"id": "mfaId",
			"code": "mfaCode"
		},
		"proposer": "proposer@example.com",
		"username": "user@example.com",
		"token": "token"
	}`
}

func TestExtractUser(t *testing.T) {
	body := makeJsonBody()
	user, err := ExtractUser(body)

	if user.Email != "email@example.com" || err != nil {
		t.Fatalf(`ExtractUser error, have %s %v, want %s nil`, user.Email, err, "email@example.com")
	}
}

func TestExtractAdminUser(t *testing.T) {
	body := makeJsonBody()
	user, err := ExtractAdminUser(body)

	if user.Email != "email@example.com" || err != nil {
		t.Fatalf(`ExtractAdminUser error, have %s %v, want %s nil`, user.Email, err, "email@example.com")
	}
}

func TestExtractGuestUser(t *testing.T) {
	body := makeJsonBody()
	user, err := ExtractGuestUser(body)

	if user.Email != "email@example.com" || err != nil {
		t.Fatalf(`ExtractGuestUser error, have %s %v, want %s nil`, user.Email, err, "email@example.com")
	}
}

func TestExtractInvite(t *testing.T) {
	body := makeJsonBody()
	user, err := ExtractInvite(body)

	if user.Proposer != "proposer@example.com" || err != nil {
		t.Fatalf(`ExtractInvite error, have %s %v, want %s nil`, user.Proposer, err, "proposer@example.com")
	}
}

func TestExtractAcceptInvite(t *testing.T) {
	body := `{"inviteeEmail": "invitee@example.com", "inviterEmail": "inviter@example.com"}`
	user, err := ExtractAcceptInvite(body)

	if user.Invitee != "invitee@example.com" || user.Inviter != "inviter@example.com" || err != nil {
		t.Fatalf(`ExtractAcceptInvite error, have %s %s %v, want %s %s nil`, user.Invitee, user.Inviter, err, "invitee@example.com", "inviter@example.com")
	}
}

func TestExtractReauth(t *testing.T) {
	body := makeJsonBody()
	user, err := ExtractReauth(body)

	if user.Email != "email@example.com" || err != nil {
		t.Fatalf(`ExtractReauth error, have %s %v, want %s nil`, user.Email, err, "email@example.com")
	}
}
