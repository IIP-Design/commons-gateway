package messages

import (
	"bytes"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type Response events.APIGatewayProxyResponse

// MarshalBody accepts any value and converts it into a stringified data object.
// TODO: properly marshall arrays of objects.
func MarshalBody(data any) ([]byte, error) {
	body, err := json.Marshal(map[string]interface{}{
		"data": data,
	})

	return body, err
}

// PrepareResponse accepts any string as an input and sets it to the message property
// of the the response body (unless there is an error when marshalling the JSON).
func PrepareResponse(body []byte) (Response, error) {
	var buf bytes.Buffer

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
