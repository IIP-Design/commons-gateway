package messages

import (
	"bytes"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type Response events.APIGatewayProxyResponse

// MarshalBody accepts any value and converts it into a stringified data object.
func MarshalBody(data any) ([]byte, error) {
	var body []byte
	var err error

	body, err = json.Marshal(map[string]any{
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
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
			"Access-Control-Allow-Origin":  "*",
			"Content-Type":                 "application/json",
		},
	}

	return resp, nil
}

// SendSuccessMessage returns a simple 200 response with a success message.
func SendSuccessMessage() (Response, error) {
	body, _ := json.Marshal(map[string]string{
		"message": "success",
	})

	return PrepareResponse(body)
}

// SendServerError accepts an error and returns it as an API Gateway response with
// a status code of 500.
func SendServerError(err error) (Response, error) {
	return Response{StatusCode: 500}, err
}
