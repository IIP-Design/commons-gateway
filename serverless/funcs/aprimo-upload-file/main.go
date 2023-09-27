package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/IIP-Design/commons-gateway/utils/aprimo"
	"github.com/IIP-Design/commons-gateway/utils/logs"
)

func uploadAprimoFile(ctx context.Context, event events.SQSEvent) error {
	var err error

	// Retrieve Aprimo auth token.
	token, err := aprimo.GetAuthToken()

	if err != nil {
		logs.LogError(err, "Unable to Authenticate Error")
		return err
	}

	for _, message := range event.Records {
		key := message.Body
		bucket := os.Getenv("SOURCE_BUCKET")

		sdkConfig, err := config.LoadDefaultConfig(ctx)

		if err != nil {
			logs.LogError(err, "Error Loading AWS Config")
			return err
		}

		s3Client := s3.NewFromConfig(sdkConfig)

		downloader := manager.NewDownloader(s3Client, func(d *manager.Downloader) {
			d.PartSize = 10 * 1024 * 1024 // 10MB per part
		})

		buffer := manager.NewWriteAtBuffer([]byte{})

		_, err = downloader.Download(context.TODO(), buffer, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})

		if err != nil {
			logs.LogError(err, "Error Retrieving S3 Object")
			return err
		}

	}

	return err
}

func main() {
	lambda.Start(uploadAprimoFile)
}
