package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	seed "github.com/IIP-Design/commons-gateway/utils/data/init"
	"github.com/IIP-Design/commons-gateway/utils/logs"
)

// seedDatabaseHandler reads a newly uploaded file from S3
func seedDatabaseHandler(ctx context.Context, event events.S3Event) error {
	var err error

	sdkConfig, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		logs.LogError(err, "Error Loading AWS Config")
		return err
	}

	s3Client := s3.NewFromConfig(sdkConfig)

	for _, record := range event.Records {
		var partMiBs int64 = 10

		downloader := manager.NewDownloader(s3Client, func(d *manager.Downloader) {
			d.PartSize = partMiBs * 1024 * 1024 // 10MB per part
		})

		buffer := manager.NewWriteAtBuffer([]byte{})

		bucket := record.S3.Bucket.Name
		key := record.S3.Object.URLDecodedKey

		_, err := downloader.Download(context.TODO(), buffer, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})

		if err != nil {
			logs.LogError(err, "Error Retrieving S3 Object")
			return err
		}

		file := bytes.NewReader(buffer.Bytes())

		csvReader := csv.NewReader(file)
		for {
			rec, err := csvReader.Read()

			if err == io.EOF {
				break
			}

			if err != nil {
				logs.LogError(err, "Error Reading CSV")
				return err
			}

			err = seed.SeedDbRecord(rec)
			if err != nil {
				logs.LogError(err, "Error Seeding DB")
				return err
			}
		}
	}

	fmt.Println("Successfully seeded the database")
	return err
}

func main() {
	lambda.Start(seedDatabaseHandler)
}
