package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"time"

	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
	"github.com/IIP-Design/commons-gateway/utils/randstr"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	LifetimeSecs = 300
)

func presignedUrlHandler(ctx context.Context, event events.APIGatewayProxyRequest) (msgs.Response, error) {
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

	key, err := randstr.RandStringBytes(24)

	if err != nil {
		logs.LogError(err, "key generation error")
		return msgs.SendServerError(err)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
	)
	if err != nil {
		logs.LogError(err, "session creation error")
		return msgs.SendServerError(err)
	}

	svc := s3.NewFromConfig(cfg)
	presigner := s3.NewPresignClient(svc)

	req, err := presigner.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s3Bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(LifetimeSecs * int64(time.Second))
	})
	if err != nil {
		logs.LogError(err, "presigning error")
		return msgs.SendServerError(err)
	}

	data := map[string]any{
		"uploadURL": req.URL,
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
	lambda.Start(presignedUrlHandler)
}
