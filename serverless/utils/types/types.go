package types

import (
	"database/sql"
	"time"
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
