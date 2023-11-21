package testHelpers

import (
	"encoding/json"
	"time"
)

func deserializeBodyData(body string) (map[string]any, error) {
	var parsed map[string]any

	b := []byte(body)
	err := json.Unmarshal(b, &parsed)

	return parsed, err
}

func DeserializeBodyArray(body string) ([]interface{}, error) {
	parsed, err := deserializeBodyData(body)
	if err != nil {
		return nil, err
	}

	return parsed["data"].([]interface{}), err
}

func FarFutureDateStr() string {
	return FarFutureDate().Format(time.RFC3339)
}

func FarFutureDate() time.Time {
	return time.Now().Add(time.Hour * 24 * 365)
}
