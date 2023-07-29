package main

import (
	"bytes"
	"encoding/json"
)

// prepareResponse accepts any string as an input and sets it to the message property
// of the the response body (unless there is an error when marshalling the JSON).
func prepareResponse(msg string) (Response, error) {
	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{
		"message": msg,
	})

	if err != nil {
		return Response{StatusCode: 404}, err
	}

	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}
