package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
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

type FileRecord struct {
	Filename string `json:"filename"`
	FileType string `json:"filetype"`
}

func ParseEventBody(body string) (FileRecord, error) {
	var parsed FileRecord

	b := []byte(body)
	err := json.Unmarshal(b, &parsed)

	return parsed, err
}

func uploadAprimoFile(ctx context.Context, event events.SQSEvent) error {
	var err error

	// Retrieve Aprimo auth token
	token, err := aprimo.GetAuthToken()

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
		fileInfo, err := ParseEventBody(message.Body)
		if err != nil {
			logs.LogError(err, "Failed to Unmarshal Body")
			return err
		}

		uri := aprimo.InitFileUpload(fileInfo.Filename, token)

		segment := 0
		readyToCommit := false

		for !readyToCommit {
			data := manager.NewWriteAtBuffer([]byte{})
			bytesDownloaded, err := downloader.Download(context.TODO(), data, &s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(fileInfo.Filename),
				Range:  aws.String(fmt.Sprintf("bytes=%d-%d", PartSize*segment, PartSize*(segment+1))),
			})

			if err != nil {
				logs.LogError(err, "Error Retrieving S3 Object")
				return err
			}

			// Send to Aprimo
			success, err := aprimo.UploadSegment(fileInfo.Filename, uri, &aprimo.FileSegment{
				Segment:  segment,
				FileType: fileInfo.FileType,
				Data:     bytes.NewBuffer(data.Bytes()),
			}, token)

			if err != nil {
				logs.LogError(err, "Aprimo Segment Upload Error")
				break
			} else if !success {
				break
			}

			segment += 1
			readyToCommit = (bytesDownloaded < PartSize)
		}

		// Commit to Aprimo
		if readyToCommit {
			ret := aprimo.CommitFileUpload(fileInfo.Filename, segment-1, uri, token)
			log.Println(ret)
		}
	}

	return err
}

func main() {
	lambda.Start(uploadAprimoFile)
}
