package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"math/big"
	"net/url"
	"os"
	"time"

	"github.com/IIP-Design/commons-gateway/utils/logs"
	msgs "github.com/IIP-Design/commons-gateway/utils/messages"
	"github.com/IIP-Design/commons-gateway/utils/security/jwt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	LetterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	LifetimeSecs = 300
)

func RandStringBytes(n int) (string, error) {
	maxVal := big.NewInt(int64(len(LetterBytes)))
	b := make([]byte, n)
	for i := range b {
		val, err := rand.Int(rand.Reader, maxVal)
		if err != nil {
			return "", err
		}
		b[i] = LetterBytes[val.Int64()]
	}

	return string(b), nil
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

	key, err := RandStringBytes(24)
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
	lambda.Start(PresignedUrlHandler)
}
