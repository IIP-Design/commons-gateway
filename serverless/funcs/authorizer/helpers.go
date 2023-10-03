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
	}
	return ""
}

// If an unknown endpoint is provided, the least permissive scope of
// `super admin` is returned.
func retrieveScopes(endpoint string) []string {
	switch endpoint {
	case "guest":
		return []string{"super admin", "admin", "guest admin", "admin"}
	case "guests":
		return []string{"super admin", "admin"}
	default:
		return []string{"super admin"}
	}
}
