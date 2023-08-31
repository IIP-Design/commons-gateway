package data

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/IIP-Design/commons-gateway/utils/logs"
)

// RequestBodyOptions represents the possible properties on the body
// JSON object sent to the serverless functions by the API Gateway.
type RequestBodyOptions struct {
	Action  string `json:"action"`
	Active  bool   `json:"active"`
	Email   string `json:"email"`
	Expires string `json:"expiration"`
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
	TeamId    string `json:"team"`
	TeamName  string `json:"teamName"`
	Username  string `json:"username"`
}

// User represents the properties required to record a user.
type User struct {
	Email     string
	NameFirst string
	NameLast  string
	Team      string
}

type AdminUser struct {
	Active bool
	User
}

type GuestUser struct {
	Expires string
	User
}

type Team struct {
	Id     string
	Name   string
	Active bool
}

// User represents the properties required to record an invite.
type Invite struct {
	Invitee User
	Inviter string
	Expires time.Time
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
	var user User

	parsed, err := ParseBodyData(body)

	email := parsed.Email
	firstName := parsed.NameFirst
	lastName := parsed.NameLast
	team := parsed.TeamId

	if err != nil {
		return user, err
	} else if email == "" || firstName == "" || lastName == "" || team == "" {
		return user, errors.New("data missing from request")
	}

	user.Email = email
	user.NameFirst = firstName
	user.NameLast = lastName
	user.Team = team

	return user, err
}

// ExtractGuestUser parses an API Gateway request body returning the data
// need to modify an existing guest user.
func ExtractGuestUser(body string) (GuestUser, error) {
	var guest GuestUser

	userData, err := ExtractUser(body)

	if err != nil {
		return guest, err
	}

	guest.Email = userData.Email
	guest.NameFirst = userData.NameFirst
	guest.NameLast = userData.NameLast
	guest.Team = userData.Team

	parsed, err := ParseBodyData(body)

	expiration := parsed.Expires

	if err != nil {
		return guest, err
	} else if expiration == "" {
		return guest, errors.New("expiration data missing from request")
	}

	guest.Expires = expiration

	return guest, err
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
	expires := parsed.Expires

	if err != nil {
		return invite, err
	} else if admin == "" || guest == "" || lastName == "" || firstName == "" || team == "" || expires == "" {
		return invite, errors.New("data missing from request")
	}

	parsedTime, err := time.Parse(time.RFC3339, expires)

	if err != nil {
		return invite, err
	}

	invite.Inviter = admin
	invite.Invitee.Email = guest
	invite.Invitee.NameFirst = firstName
	invite.Invitee.NameLast = lastName
	invite.Invitee.Team = team
	invite.Expires = parsedTime

	return invite, err
}
