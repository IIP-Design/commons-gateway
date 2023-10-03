package main

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	jwtv5 "github.com/golang-jwt/jwt/v5"

	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/IIP-Design/commons-gateway/utils/security/jwt"
)

type ARNInfo struct {
	AccountId string
	APIId     string
	Method    string
	Resource  string
	Region    string
	Stage     string
}

// parseMethodArn breaks up the information provided by the API
// Gateway endpoint into usable chunks.
func parseMethodArn(arn string) ARNInfo {
	parts := strings.Split(arn, ":")

	var arnInfo ARNInfo

	arnInfo.Region = parts[3]
	arnInfo.AccountId = parts[4]

	// Parse the gateway endpoints
	apiGatewayPath := strings.Split(parts[5], "/")
	arnInfo.APIId = apiGatewayPath[0]
	arnInfo.Stage = apiGatewayPath[1]
	arnInfo.Method = apiGatewayPath[2]

	segmentCount := len(apiGatewayPath)
	tail := apiGatewayPath[3:segmentCount]

	arnInfo.Resource = strings.Join(tail, "/")

	return arnInfo
}

// setPolicyStatement constructs a lambda execution ARN based on
// the information found in the API Gateway method ARN.
func setPolicyStatement(effect Effect, arnInfo ARNInfo) events.IAMPolicyStatement {
	resourceArn := "arn:aws:execute-api:" +
		arnInfo.Region + ":" +
		arnInfo.AccountId + ":" +
		arnInfo.APIId + "/" +
		arnInfo.Stage + "/" +
		arnInfo.Method + "/" +
		arnInfo.Resource

	statement := events.IAMPolicyStatement{
		Effect:   effect.String(),
		Action:   []string{"execute-api:Invoke"},
		Resource: []string{resourceArn},
	}

	return statement
}

// rejectRequest outright rejects a request sending back an 401 Unauthorized error.
func rejectRequest(status int) (events.APIGatewayCustomAuthorizerResponse, error) {
	var msg string

	switch status {
	case 403:
		msg = "Forbidden"
	case 401:
		msg = "Unauthorized"
	default:
		msg = "Unauthorized"
	}

	return events.APIGatewayCustomAuthorizerResponse{}, errors.New(msg)
}

// handleAuthorizationRequest orchestrates the verification of the provided
// authorization token and grants subsequent access to invoke the Lambda.
func handleAuthorizationRequest(
	ctx context.Context,
	event events.APIGatewayCustomAuthorizerRequest,
) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := event.AuthorizationToken
	arnInfo := parseMethodArn(event.MethodArn)

	// Short circuit if no token provided.
	if token == "" {
		logs.LogError(errors.New("missing authorization token"), "Authorization Token Error")
		return rejectRequest(401)
	}

	// Verify the token is valid.
	err := jwt.CheckAuthToken(token, retrieveScopes(arnInfo.Resource))

	if err != nil {
		logs.LogError(err, "Error Validating JWT")

		if errors.Is(err, jwtv5.ErrTokenExpired) {
			return rejectRequest(401)
		} else {
			return rejectRequest(403)
		}
	}

	// Construct the IAM policy allowing the user to invoke the Lambda.
	statement := setPolicyStatement(Allow, arnInfo)

	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: "",
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version:   "2012-10-17",
			Statement: []events.IAMPolicyStatement{statement},
		},
	}, nil
}

func main() {
	lambda.Start(handleAuthorizationRequest)
}
