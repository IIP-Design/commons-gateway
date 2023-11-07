package data

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/IIP-Design/commons-gateway/utils/types"
)

// ParseBodyData converts the serialized JSON string provided in the body
// of the API Gateway request into a usable data format.
func ParseBodyData(body string) (types.RequestBodyOptions, error) {
	var parsed types.RequestBodyOptions

	b := []byte(body)
	err := json.Unmarshal(b, &parsed)

	if err != nil {
		logs.LogError(err, "Failed to Unmarshal Body")
	}

	return parsed, err
}

// ExtractUser parses an API Gateway request body returning the data
// need to create a new user.
func ExtractUser(body string) (types.User, error) {
	var user types.User

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
func ExtractAdminUser(body string) (types.AdminUser, error) {
	var admin types.AdminUser

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
func ExtractGuestUser(body string) (types.GuestUser, error) {
	var guest types.GuestUser

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
func ExtractInvite(body string) (types.Invite, error) {
	var invite types.Invite

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

func ExtractAcceptInvite(body string) (types.AcceptInvite, error) {
	var invite types.AcceptInvite

	b := []byte(body)
	err := json.Unmarshal(b, &invite)

	if err != nil {
		logs.LogError(err, "Failed to Unmarshal Invite")
	}

	return invite, err
}

func ExtractReauth(body string) (types.GuestReauth, error) {
	var ret types.GuestReauth

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
