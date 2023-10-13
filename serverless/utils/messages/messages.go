package messages

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

type Response events.APIGatewayProxyResponse

func marshalResponse(data any, prop string) ([]byte, error) {
	var body []byte
	var err error

	body, err = json.Marshal(map[string]any{
		prop: data,
	})

	return body, err
}

// MarshalBody accepts any value and converts it into a stringified data object.
func MarshalBody(data any) ([]byte, error) {
	return marshalResponse(data, "data")
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
	log.Println(err.Error())
	return Response{
		StatusCode:      500,
		IsBase64Encoded: false,
		Body:            err.Error(),
		Headers: map[string]string{
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
			"Access-Control-Allow-Origin":  "*",
			"Content-Type":                 "application/json",
		},
	}, nil
}

// statusCodeToBody returns standardized error messages
// based on the provided status code
func statusCodeToBody(statusCode int) string {
	var code string

	switch statusCode {
	case 401:
		code = "unauthorized"
	case 403:
		code = "forbidden"
	case 500:
	default:
		code = "internal error"
	}

	return code
}

func SendAuthError(err error, statusCode int) (Response, error) {
	var msg string

	if err != nil {
		msg = err.Error()
	} else {
		msg = statusCodeToBody(statusCode)
	}

	body, _ := marshalResponse(msg, "error")

	var buf bytes.Buffer

	json.HTMLEscape(&buf, body)

	headers := map[string]string{
		"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
		"Access-Control-Allow-Origin":  "*",
		"Content-Type":                 "application/json",
	}

	// Set a retry after header if the response is Too Many Requests.
	if statusCode == 429 {
		headers["Retry-After"] = "900"
	}

	resp := Response{
		StatusCode:      statusCode,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers:         headers,
	}
	return resp, nil
}
