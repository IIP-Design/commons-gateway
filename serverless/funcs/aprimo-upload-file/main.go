package main

import (
	"context"
	"fmt"
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

const (
	PartSize = 10 * 1024 * 1024 // 10MB per part
)

func uploadAprimoFile(ctx context.Context, event events.SQSEvent) error {
	var err error

	// Retrieve Aprimo auth token.  FIXME: Add back token later
	_, err = aprimo.GetAuthToken()

	if err != nil {
		logs.LogError(err, "Unable to Authenticate Error")
		return err
	}

	sdkConfig, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		logs.LogError(err, "Error Loading AWS Config")
		return err
	}

	s3Client := s3.NewFromConfig(sdkConfig)
	bucket := os.Getenv("SOURCE_BUCKET")

	downloader := manager.NewDownloader(s3Client, func(d *manager.Downloader) {
		d.PartSize = PartSize
	})

	for _, message := range event.Records {
		key := message.Body

		segment := 0
		readyToCommit := false

		for !readyToCommit {
			buffer := manager.NewWriteAtBuffer([]byte{})
			bytesDownloaded, err := downloader.Download(context.TODO(), buffer, &s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
				Range:  aws.String(fmt.Sprintf("bytes=%d-%d", PartSize*segment, PartSize*(segment+1))),
			})

			if err != nil {
				logs.LogError(err, "Error Retrieving S3 Object")
				return err
			}

			// Send to Aprimo

			segment += 1
			readyToCommit = (bytesDownloaded < PartSize)
		}

		// Commit to Aprimo
	}

	return err
}

func main() {
	lambda.Start(uploadAprimoFile)
}
