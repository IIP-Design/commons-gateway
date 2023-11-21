package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestPresignedUrl(t *testing.T) {
	event := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"contentType": "image/png",
			"fileName":    "image.png",
		},
	}

	resp, err := presignedUrlHandler(context.TODO(), event)

	if resp.StatusCode != 200 || err != nil {
		t.Fatalf(`presignedUrlHandler failed, have status %d and %v, want 200 and nil`, resp.StatusCode, err)
	}
}
