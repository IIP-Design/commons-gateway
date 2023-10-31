package data

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/IIP-Design/commons-gateway/utils/logs"
)

// UserBodyOptions represents the possible properties on the body
// JSON object sent to the serverless functions by the API Gateway
// for operations that update guest or admin users.
type UserBodyOptions struct {
	Email      string `json:"email"`
	Expires    string `json:"expiration"`
	NameFirst  string `json:"givenName"`
	NameLast   string `json:"familyName"`
	Role       string `json:"role"`
	TeamId     string `json:"team"`
	TeamName   string `json:"teamName"`
	AprimoName string `json:"teamAprimo"`
}

type MFARequest struct {
	Id   string `json:"id"`
	Code string `json:"code"`
}

// RequestBodyOptions represents the possible properties on the body
// JSON object sent to the serverless functions by the API Gateway.
type RequestBodyOptions struct {
	UserBodyOptions
	Active   bool            `json:"active"`
	Hash     string          `json:"hash"`
	Invitee  UserBodyOptions `json:"invitee"`
	Inviter  string          `json:"inviter"`
	Admin    string          `json:"admin"`
	MFA      MFARequest      `json:"mfa"`
	Proposer string          `json:"proposer"`
	Username string          `json:"username"`
	Token    string          `json:"token"`
}

// User represents the properties required to record a user.
type User struct {
	Email     string `json:"email"`
	NameFirst string `json:"givenName"`
	NameLast  string `json:"familyName"`
	Role      string `json:"role"`
	Team      string `json:"team"`
}

// AdminUser extends the base User struct with unique admin properties.
type AdminUser struct {
	Active bool
	User
}

// GuestUser extends the base User struct with unique guest properties.
type GuestUser struct {
	Expires string `json:"expires"`
	Pending bool   `json:"pending"`
	User
}

// GuestInvite extends the GuestUser struct with unique invite properties.
type GuestInvite struct {
	DateInvited string
	Proposer    string
	GuestUser
}

type UploaderUser struct {
	GuestUser
	DateInvited string
	Proposer    sql.NullString
	Inviter     sql.NullString
	Pending     bool
}

// Team represents the properties required to record a team.
type Team struct {
	Id         string
	Name       string
	AprimoName string
	Active     bool
}

// User represents the properties required to record an invite.
type Invite struct {
	Invitee  User
	Inviter  string
	Proposer string
	Expires  time.Time
}

type AcceptInvite struct {
	Invitee string `json:"inviteeEmail"`
	Inviter string `json:"inviterEmail"`
}

type GuestUnlockInitEvent struct {
	Username string `json:"username"`
}

type GuestReauth struct {
	Email   string    `json:"email"`
	Admin   string    `json:"admin"`
	Expires time.Time `json:"expiration"`
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
	role := parsed.Role
	team := parsed.TeamId

	if err != nil {
		return user, err
	} else if email == "" || firstName == "" || lastName == "" || role == "" || team == "" {
		return user, errors.New("data missing from request")
	}

	user.Email = email
	user.NameFirst = firstName
	user.NameLast = lastName
	user.Role = role
	user.Team = team

	return user, err
}

// ExtractAdminUser parses an API Gateway request body returning the data
// need to modify an existing admin user.
func ExtractAdminUser(body string) (AdminUser, error) {
	var admin AdminUser

	userData, err := ExtractUser(body)

	if err != nil {
		return admin, err
	}

	admin.Email = userData.Email
	admin.NameFirst = userData.NameFirst
	admin.NameLast = userData.NameLast
	admin.Role = userData.Role
	admin.Team = userData.Team

	parsed, err := ParseBodyData(body)

	active := parsed.Active

	if err != nil {
		return admin, err
	}

	admin.Active = active

	return admin, err
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
	guest.Role = userData.Role
	guest.Team = userData.Team

	return guest, err
}

// ExtractInvite parses an API Gateway request body returning the data
// need to create a guest user invitation.
func ExtractInvite(body string) (Invite, error) {
	var invite Invite

	parsed, err := ParseBodyData(body)

	admin := parsed.Inviter
	proposer := parsed.Proposer
	guest := parsed.Invitee.Email
	firstName := parsed.Invitee.NameFirst
	lastName := parsed.Invitee.NameLast
	role := parsed.Invitee.Role
	team := parsed.Invitee.TeamId
	expires := parsed.Expires

	if err != nil {
		return invite, err
	} else if guest == "" || lastName == "" || firstName == "" || team == "" || expires == "" {
		return invite, errors.New("data missing from request")
	} else if admin == "" && proposer == "" {
		return invite, errors.New("must supply admin or proposer")
	}

	// Default the role to guest if not provided.
	if role == "" {
		role = "guest"
	}

	parsedTime, err := time.Parse(time.RFC3339, expires)

	if err != nil {
		return invite, err
	}

	invite.Inviter = admin
	invite.Proposer = proposer
	invite.Invitee.Email = guest
	invite.Invitee.NameFirst = firstName
	invite.Invitee.NameLast = lastName
	invite.Invitee.Role = role
	invite.Invitee.Team = team
	invite.Expires = parsedTime

	return invite, err
}

func ExtractAcceptInvite(body string) (AcceptInvite, error) {
	var invite AcceptInvite

	b := []byte(body)
	err := json.Unmarshal(b, &invite)

	if err != nil {
		logs.LogError(err, "Failed to Unmarshal Invite")
	}

	return invite, err
}

func ExtractReauth(body string) (GuestReauth, error) {
	var ret GuestReauth

	parsed, err := ParseBodyData(body)
	if err != nil {
		return ret, err
	}

	email := parsed.Email
	admin := parsed.Admin
	expires := parsed.Expires

	parsedTime, err := time.Parse(time.RFC3339, expires)

	if err != nil {
		return ret, err
	}

	ret.Email = email
	ret.Admin = admin
	ret.Expires = parsedTime

	return ret, err
}
