package testHelpers

import (
	"encoding/json"
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

func DeserializeBodyObject(body string) (map[string]any, error) {
	parsed, err := deserializeBodyData(body)
	if err != nil {
		return nil, err
	}

	return parsed["data"].(map[string]any), err
}
