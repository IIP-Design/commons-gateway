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
	"github.com/IIP-Design/commons-gateway/utils/queue"
)

const (
	PartSize = 10 * 1024 * 1024 // 10MB per part
)

func ParseEventBody(body string) (aprimo.FileRecordInitEvent, error) {
	var parsed aprimo.FileRecordInitEvent

	b := []byte(body)
	err := json.Unmarshal(b, &parsed)

	return parsed, err
}

func SendUpdateEvent(aprimoId string, filename string, fileToken string) (string, error) {
	var messageId string
	var err error

	event := aprimo.FileRecordUpdateEvent{
		AprimoId:  aprimoId,
		Filename:  filename,
		FileToken: fileToken,
	}

	json, err := json.Marshal(event)

	if err != nil {
		logs.LogError(err, "Failed to Marshal SQS Body")
		return messageId, err
	}

	queueUrl := os.Getenv("RECORD_UPDATE_QUEUE")

	// Send the message to SQS.
	return queue.SendToQueue(string(json), queueUrl)
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
			uploadToken, err := aprimo.CommitFileUpload(fileInfo.Filename, segment-1, uri, token)
			if err != nil {
				logs.LogError(err, "Aprimo File Commit Error")
			} else {
				log.Println(uploadToken)
				messageId, err := SendUpdateEvent(fileInfo.AprimoId, fileInfo.Filename, uploadToken)
				if err != nil {
					logs.LogError(err, "send record update event error")
				} else {
					log.Println(messageId)
				}

			}
		} else {
			log.Println("Not ready to commit")
		}
	}

	return err
}

func main() {
	lambda.Start(uploadAprimoFile)
}
