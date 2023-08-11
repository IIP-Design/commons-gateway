package data

import (
	"encoding/json"

	"github.com/IIP-Design/commons-gateway/utils/logs"
)

// RequestBodyOptions represents the possible properties on the body
// JSON object sent to the serverless functions by the API Gateway.
type RequestBodyOptions struct {
	Action   string `json:"action"`
	Email    string `json:"email"`
	Hash     string `json:"hash"`
	Invitee  string `json:"invitee"`
	Inviter  string `json:"inviter"`
	Username string `json:"username"`
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
