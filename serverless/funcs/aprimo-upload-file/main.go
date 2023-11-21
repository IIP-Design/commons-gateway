package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

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

type WrappedS3Events struct {
	Records []events.S3EventRecord `json:"Records"`
}

const (
	PartSize        = 19 * 1024 * 1024 // 19 MB per part
	PartsToDownload = 10
	S3DownloadBytes = PartSize * PartsToDownload
)

// Needed before go 1.21
func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// lookupFileType returns the file type info if file has not already been uploaded
// to Aprimo. Duplication is not an error, but is a reason to skip re-processing
func lookupFileType(key string) (string, error) {
	pool := data.ConnectToDB()
	defer pool.Close()

	var fileType string
	var aprimoUploadToken sql.NullString

	query := "SELECT file_type, aprimo_upload_token FROM uploads WHERE s3_id = $1"
	err := pool.QueryRow(query, key).Scan(&fileType, &aprimoUploadToken)

	// There is a value for the token, so it's already been uploaded
	if aprimoUploadToken.Valid {
		return "", err
	} else {
		return fileType, err
	}
}

// If the file has been transferred to Aprimo, record the token for (1) possible retry and (2) deduplication
func markFileUpload(key string, uploadToken string) error {
	pool := data.ConnectToDB()
	defer pool.Close()

	query := "UPDATE uploads SET aprimo_upload_token = $1, aprimo_upload_dt = NOW() WHERE s3_id = $2"
	_, err := pool.Exec(query, uploadToken, key)

	return err
}

func uploadSmallFile(key string, token string, downloader *manager.Downloader, bucket string, fileType string) (string, error) {
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

func uploadFileInSegments(key string, token string, downloader *manager.Downloader, bucket string, fileType string, fileSize int64) (string, error) {
	var uploadToken string
	var err error

	uri := aprimo.InitFileUpload(key, token)

	buf := make([]byte, S3DownloadBytes)
	segment := 0
	downloadBlock := 0
	readyToCommit := false

	var t1 int64
	var t2 int64

	for !readyToCommit {
		// NB: Range appears to be inclusive at both ends
		s3DownloadStartByte := S3DownloadBytes * downloadBlock
		s3DownloadEndByte := min(int64(S3DownloadBytes*(downloadBlock+1)-1), fileSize)

		// DBG
		fmt.Printf("bytes=%d-%d\n", s3DownloadStartByte, s3DownloadEndByte)

		data := manager.NewWriteAtBuffer(buf)

		t1 = time.Now().UnixMilli()
		bytesDownloaded, err := downloader.Download(context.TODO(), data, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Range:  aws.String(fmt.Sprintf("bytes=%d-%d", s3DownloadStartByte, s3DownloadEndByte)),
		})
		t2 = time.Now().UnixMilli()

		if err != nil {
			logs.LogError(err, "Error Retrieving S3 Object")
			break
		} else {
			// DBG
			fmt.Printf("Downloaded %d bytes in %d ms\n", bytesDownloaded, t2-t1)
		}

		dataBytes := data.Bytes()

		segmentsDownloaded := (bytesDownloaded / PartSize)
		lastSegmentIsShort := bytesDownloaded%PartSize > 0
		if lastSegmentIsShort {
			// DBG
			fmt.Printf("Excess bytes: %d \n", bytesDownloaded%PartSize)
			segmentsDownloaded += 1
		}

		// DBG
		fmt.Printf("Downloaded %d segments\n", segmentsDownloaded)

		for seg := int64(0); seg < segmentsDownloaded; seg++ {
			start := PartSize * int64(seg)
			end := min(int64(PartSize*(seg+1)), bytesDownloaded)

			// Send to Aprimo
			t1 = time.Now().UnixMilli()
			success, err := aprimo.UploadSegment(key, uri, &aprimo.FileSegment{
				Segment:  segment,
				FileType: fileType,
				Data:     bytes.NewBuffer(dataBytes[start:end]),
			}, token)
			t2 = time.Now().UnixMilli()

			if err != nil {
				logs.LogError(err, "Aprimo Segment Upload Error")
				break
			} else if !success {
				logs.LogError(err, "Failed to succeed")
				break
			}

			// DBG
			fmt.Printf("Uploaded bytes %d to %d in %d ms\n", start, end, t2-t1)

			segment += 1
		}

		readyToCommit = (bytesDownloaded < S3DownloadBytes)

		downloadBlock += 1
	}

	// Commit to Aprimo
	if readyToCommit {
		// Segments are zero-indexed but we need to indicate the total number of segments
		// So the final increment of the loop should give the proper value
		uploadToken, err = aprimo.CommitFileUpload(key, segment, uri, token)
	} else {
		log.Println("Not ready to commit")
	}

	return uploadToken, err
}

func sendRecordEvent(key string, fileType string, fileToken string) (string, error) {
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

	queueUrl := os.Getenv("RECORD_CREATE_QUEUE")

	// Send the message to SQS.
	return queue.SendToQueue(string(json), queueUrl)
}

// extractS3DataFromSqsEvent retrieves the S3 event embedded in the SQS message.
func extractS3DataFromSqsEvent(record events.SQSMessage) []events.S3EventRecord {
	var parsed WrappedS3Events
	json.Unmarshal([]byte(record.Body), &parsed)

	return parsed.Records
}

func handleUploadFileToAprimo(ctx context.Context, event events.SQSEvent) error {
	var err error

	// Retrieve Aprimo auth token
	token, err := aprimo.GetAuthToken()

	if err != nil {
		logs.LogError(err, "Unable to Authenticate Error")
		return err
	}

	// Prepare AWS services
	sdkConfig, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		logs.LogError(err, "Error Loading AWS Config")
		return err
	}

	s3Client := s3.NewFromConfig(sdkConfig)

	downloader := manager.NewDownloader(s3Client, func(d *manager.Downloader) {
		d.PartSize = PartSize
	})

	// We are receiving SQS record(s) for increased durability...
	for _, e := range event.Records {
		log.Printf("Message ID: %s\n", e.MessageId)

		// ...but they are wrapping one or more (should be exactly one, but handle 1..N) S3 events...
		r := extractS3DataFromSqsEvent(e)

		// ...and need to be processed as such
		for _, record := range r {
			bucket := record.S3.Bucket.Name
			key := record.S3.Object.Key
			size := record.S3.Object.Size

			fileType, err := lookupFileType(key)

			if err != nil {
				logs.LogError(err, "failed to lookup file type")
				return err
			}

			if fileType != "" { // Have a new record
				var uploadToken string

				// Concerted upload for small files, segmented if necessary
				if size <= PartSize {
					uploadToken, err = uploadSmallFile(key, token, downloader, bucket, fileType)
				} else {
					uploadToken, err = uploadFileInSegments(key, token, downloader, bucket, fileType, size)
				}

				if err == nil {
					messageId, err := sendRecordEvent(key, fileType, uploadToken)
					if err != nil {
						logs.LogError(err, "send record event error")
						return err
					}
					log.Printf("Object %s sent onwards with message ID %s\n", key, messageId)

					err = markFileUpload(key, uploadToken)
					if err != nil {
						logs.LogError(err, "mark file upload error")
						return err
					}
				} else {
					logs.LogError(err, "failed to upload file")
					return err
				}
			} else { // Not a new record
				log.Printf("Object %s has already been uploaded, but the event was not deleted\n", key)
			}
		}
	}

	return nil
}

func main() {
	lambda.Start(handleUploadFileToAprimo)
}
