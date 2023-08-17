package data

import (
	"encoding/json"
	"errors"

	"github.com/IIP-Design/commons-gateway/utils/logs"
)

// RequestBodyOptions represents the possible properties on the body
// JSON object sent to the serverless functions by the API Gateway.
type RequestBodyOptions struct {
	Action  string `json:"action"`
	Email   string `json:"email"`
	Hash    string `json:"hash"`
	Invitee struct {
		Email     string `json:"email"`
		NameFirst string `json:"givenName"`
		NameLast  string `json:"familyName"`
		Team      string `json:"team"`
	} `json:"invitee"`
	Inviter   string `json:"inviter"`
	NameFirst string `json:"givenName"`
	NameLast  string `json:"familyName"`
	Team      string `json:"team"`
	Username  string `json:"username"`
}

// User represents the properties required to record a user.
type User struct {
	Email     string
	NameFirst string
	NameLast  string
	Team      string
}

// User represents the properties required to record an invite.
type Invite struct {
	Invitee User
	Inviter string
}

// ParseBodyData converts the serialized JSON string provided in the body
// of the API Gateway request into a usable data format.
func ParseBodyData(body string) (RequestBodyOptions, error) {
	var parsed RequestBodyOptions

	b := []byte(body)
	err := json.Unmarshal(b, &parsed)

	if err != nil {
		logs.LogError(err, "Failed to Unmarshal Body")
	}

	return parsed, err
}

// ExtractUser parses an API Gateway request body returning the data
// need to create a new user.
func ExtractUser(body string) (User, error) {
	var admin User

	parsed, err := ParseBodyData(body)

	adminEmail := parsed.Email
	firstName := parsed.NameFirst
	lastName := parsed.NameLast
	team := parsed.Team

	if err != nil {
		return admin, err
	} else if adminEmail == "" || firstName == "" || lastName == "" || team == "" {
		return admin, errors.New("data missing from request")
	}

	admin.Email = adminEmail
	admin.NameFirst = firstName
	admin.NameLast = lastName
	admin.Team = team

	return admin, err
}

// ExtractInvite parses an API Gateway request body returning the data
// need to create a guest user invitation.
func ExtractInvite(body string) (Invite, error) {
	var invite Invite

	parsed, err := ParseBodyData(body)

	admin := parsed.Inviter
	guest := parsed.Invitee.Email
	firstName := parsed.Invitee.NameFirst
	lastName := parsed.Invitee.NameLast
	team := parsed.Invitee.Team

	if err != nil {
		return invite, err
	} else if admin == "" || guest == "" || lastName == "" || firstName == "" || team == "" {
		return invite, errors.New("data missing from request")
	}

	invite.Inviter = admin
	invite.Invitee.Email = guest
	invite.Invitee.NameFirst = firstName
	invite.Invitee.NameLast = lastName
	invite.Invitee.Team = team

	return invite, err
}
