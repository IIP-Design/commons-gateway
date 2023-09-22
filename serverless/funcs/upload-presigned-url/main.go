package main

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"net/url"
	"os"
	"time"

	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
	"github.com/IIP-Design/commons-gateway/utils/security/jwt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func PresignedUrlHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
	code, err := jwt.RequestIsAuthorized(event, []string{"super admin", "admin", "guest admin"})
	if err != nil {
		return msgs.SendAuthError(err, code)
	}

	rawContentType := event.QueryStringParameters["contentType"]
	if rawContentType == "" {
		return msgs.SendServerError(errors.New("no content type submitted"))
	}
	contentType, err := url.PathUnescape(rawContentType)
	if err != nil {
		logs.LogError(err, "content-type decode error")
		return msgs.SendServerError(err)
	}

	var awsRegion = os.Getenv("AWS_REGION")
	var s3Bucket = os.Getenv("S3_UPLOAD_BUCKET")

	key := RandStringBytes(24)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)
	if err != nil {
		logs.LogError(err, "session creation error")
		return msgs.SendServerError(err)
	}

	svc := s3.New(sess)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(s3Bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	})
	urlStr, err := req.Presign(300 * time.Second)

	if err != nil {
		logs.LogError(err, "presigning error")
		return msgs.SendServerError(err)
	}

	data := map[string]any{
		"uploadURL": urlStr,
		"key":       key,
	}

	body, err := json.Marshal(data)

	if err != nil {
		logs.LogError(err, "data marshalling")
		return msgs.SendServerError(err)
	}

	return msgs.PrepareResponse(body)
}

func main() {
	lambda.Start(PresignedUrlHandler)
}
