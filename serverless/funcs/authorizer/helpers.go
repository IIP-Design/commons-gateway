package main

type Effect int

const (
	Allow Effect = iota
	Deny
)

func (e Effect) String() string {
	switch e {
	case Allow:
		return "Allow"
	case Deny:
		return "Deny"
	default:
		return ""
	}
}

type Scopes int

const (
	All         Scopes = iota // all users
	AllAdmins                 // super admins, admins, external admins
	GuestAdmins               // external partner admins
	StateAdmins               // super admins, admins
	SuperAdmins               // super admins
)

func (s Scopes) Array() []string {
	switch s {
	case All:
		return []string{"super admin", "admin", "guest admin", "guest"}
	case AllAdmins:
		return []string{"super admin", "admin", "guest admin"}
	case GuestAdmins:
		return []string{"guest admin"}
	case StateAdmins:
		return []string{"super admin", "admin"}
	case SuperAdmins:
		return []string{"super admin"}
	default:
		return []string{}
	}
}

// If an unknown endpoint is provided, the least permissive scope of
// `super admin` is returned.
func retrieveScopes(endpoint string, method string) []string {
	switch endpoint {
	case "admin":
		if method == "GET" {
			return StateAdmins.Array()
		} else if method == "DELETE" {
			return SuperAdmins.Array()
		} else if method == "POST" {
			return SuperAdmins.Array()
		} else if method == "PUT" {
			return SuperAdmins.Array()
		}
	case "admins":
		return SuperAdmins.Array()
	case "creds/propose":
		return GuestAdmins.Array()
	case "creds/provision":
		return StateAdmins.Array()
	case "guest":
		if method == "GET" {
			return All.Array()
		} else if method == "DELETE" {
			return AllAdmins.Array()
		} else if method == "PUT" {
			return AllAdmins.Array()
		}
	case "guest/approve":
		return StateAdmins.Array()
	case "guest/reauth":
		return AllAdmins.Array()
	case "guests":
		return StateAdmins.Array()
	case "guests/pending":
		return StateAdmins.Array()
	case "guests/uploaders":
		return GuestAdmins.Array()
	case "team":
		if method == "POST" {
			return SuperAdmins.Array()
		} else if method == "PUT" {
			return SuperAdmins.Array()
		}
	case "teams":
		return AllAdmins.Array()
	case "upload":
		return All.Array()
	default:
		return SuperAdmins.Array()
	}

	return SuperAdmins.Array()
}
