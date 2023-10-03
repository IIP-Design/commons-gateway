package main

import (
	"context"
	"encoding/json"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RequestBody struct {
	Sql string `json:"sql"`
}

func parseRequest(body string) (RequestBody, error) {
	var parsed RequestBody

	b := []byte(body)
	err := json.Unmarshal(b, &parsed)

	return parsed, err
}

func execQuery(ctx context.Context, event events.APIGatewayProxyRequest) error {
	requestBody, err := parseRequest(event.Body)
	if err != nil {
		logs.LogError(err, "Failed to Unmarshal Body")
		return err
	}

	pool := data.ConnectToDB()
	defer pool.Close()

	_, err = pool.Exec(requestBody.Sql)

	return err
}

func main() {
	lambda.Start(execQuery)
}
