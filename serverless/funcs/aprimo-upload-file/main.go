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
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/IIP-Design/commons-gateway/utils/queue"
)

const (
	PartSize = 15 * 1024 * 1024 // 15MB per part
)

func LookupFileType(key string) (string, error) {
	pool := data.ConnectToDB()
	defer pool.Close()

	var fileType string

	query := "SELECT file_type FROM uploads WHERE s3_id = $1"
	err := pool.QueryRow(query, key).Scan(&fileType)

	return fileType, err
}

func UploadSmallFile(key string, token string, downloader *manager.Downloader, bucket string, fileType string) (string, error) {
	data := manager.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(context.TODO(), data, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		logs.LogError(err, "Error Retrieving S3 Object")
		return "", err
	}

	// Send to Aprimo
	return aprimo.UploadFile(key, fileType, bytes.NewBuffer(data.Bytes()), token)
}

func UploadFileSegments(key string, token string, downloader *manager.Downloader, bucket string, fileType string) (string, error) {
	var uploadToken string
	var err error

	uri := aprimo.InitFileUpload(key, token)

	segment := 0
	readyToCommit := false

	for !readyToCommit {
		data := manager.NewWriteAtBuffer([]byte{})
		bytesDownloaded, err := downloader.Download(context.TODO(), data, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Range:  aws.String(fmt.Sprintf("bytes=%d-%d", PartSize*segment, PartSize*(segment+1))),
		})

		if err != nil {
			logs.LogError(err, "Error Retrieving S3 Object")
			break
		}

		// Send to Aprimo
		success, err := aprimo.UploadSegment(key, uri, &aprimo.FileSegment{
			Segment:  segment,
			FileType: fileType,
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
		uploadToken, err = aprimo.CommitFileUpload(key, segment-1, uri, token)
	} else {
		log.Println("Not ready to commit")
	}

	return uploadToken, err
}

func SendRecordEvent(key string, fileType string, fileToken string) (string, error) {
	var messageId string
	var err error

	event := aprimo.FileRecordInitEvent{
		Key:       key,
		FileType:  fileType,
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

func uploadAprimoFile(ctx context.Context, event events.S3Event) error {
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

	downloader := manager.NewDownloader(s3Client, func(d *manager.Downloader) {
		d.PartSize = PartSize
	})

	for _, record := range event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.Key
		size := record.S3.Object.Size

		fileType, err := LookupFileType(key)
		if err != nil {
			logs.LogError(err, "failed to lookup file type")
			return err
		}

		var uploadToken string

		if size <= PartSize {
			uploadToken, err = UploadSmallFile(key, token, downloader, bucket, fileType)
		} else {
			uploadToken, err = UploadFileSegments(key, token, downloader, bucket, fileType)
		}

		if err == nil {
			log.Println(uploadToken)
			messageId, err := SendRecordEvent(key, fileType, uploadToken)
			if err != nil {
				logs.LogError(err, "send record event error")
			} else {
				log.Println(messageId)
			}
		}
	}

	return err
}

func main() {
	lambda.Start(uploadAprimoFile)
}
