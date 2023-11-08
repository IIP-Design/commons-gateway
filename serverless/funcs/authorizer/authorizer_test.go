package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

const (
	TEST_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2OTk0NjAyMzUsImZpcnN0TG9naW4iOmZhbHNlLCJzY29wZSI6Imd1ZXN0IiwidXNlciI6InRlc3RAZXhhbXBsZS5jb20ifQ.zbvLuFR-HuJblu3nuwjKDKD5db_gfmb9xCFIgeZZHuU"
	TEST_ARN   = "arn:aws:execute-api:us-east-1:123456789012:abcdef123/test/POST/upload"
)

func TestParseArn(t *testing.T) {
	parsed := parseMethodArn(TEST_ARN)
	if parsed.Region != "us-east-1" {
		t.Fatalf("parseMethodArn incorrect region: %s", parsed.Region)
	} else if parsed.AccountId != "123456789012" {
		t.Fatalf("parseMethodArn incorrect account: %s", parsed.AccountId)
	} else if parsed.Stage != "test" {
		t.Fatalf("parseMethodArn incorrect stage: %s", parsed.Stage)
	}
}

func TestSetPolicyStatement(t *testing.T) {
	stmt := setPolicyStatement(Allow, parseMethodArn(TEST_ARN))
	if stmt.Effect != "Allow" {
		t.Fatalf("setPolicyStatement incorrect effect: %s", stmt.Effect)
	} else if stmt.Action[0] != "execute-api:Invoke" {
		t.Fatalf("setPolicyStatement incorrect action: %s", stmt.Action[0])
	} else if stmt.Resource[0] != TEST_ARN {
		t.Fatalf("setPolicyStatement incorrect resource: %s", stmt.Resource[0])
	}
}

func TestRejectRequest(t *testing.T) {
	_, err := rejectRequest(403)
	if err.Error() != "Forbidden" {
		t.Fatalf("rejectRequest incorrect message: %s", err.Error())
	}
}

func TestRejectRequestIncorrectValue(t *testing.T) {
	_, err := rejectRequest(422)
	if err.Error() != "Unauthorized" {
		t.Fatalf("rejectRequest incorrect message: %s", err.Error())
	}
}

func TestHandleRequest(t *testing.T) {
	event := events.APIGatewayCustomAuthorizerRequest{
		AuthorizationToken: TEST_TOKEN,
		MethodArn:          TEST_ARN,
	}

	resp, err := handleAuthorizationRequest(context.TODO(), event)

	if resp.PrincipalID != "*" || err != nil {
		t.Fatalf("handleAuthorizationRequest failure: %s %v, want \"*\" nil", resp.PrincipalID, err)
	}
}
