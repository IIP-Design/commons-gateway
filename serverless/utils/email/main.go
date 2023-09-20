package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
)

type ProvisionCredsData struct {
	Invitee     data.User `json:"invitee"`
	Inviter     data.User `json:"inviter"`
	TmpPassword string    `json:"tmpPassword"`
	Url         string    `json:"url"`
}

type RequestSupportStaffData struct {
	Proposer data.User `json:"externalTeamLead"`
	Invitee  data.User `json:"supportStaffuser"`
	Url      string    `json:"url"`
}

type TwoFactorAuthData struct {
	User data.User `json:"user"`
	Code string    `json:"verificationCode"`
}

func serialize(data any) (string, error) {
	serial, err := json.Marshal(data)

	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	json.HTMLEscape(&buf, serial)
	return buf.String(), nil
}

func getQueueUrl(svc *sqs.SQS, queueName string) (string, error) {
	queueUrlOutput, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: aws.String(queueName)})

	if err != nil {
		return "", err
	}

	queueUrl := queueUrlOutput.QueueUrl

	return *queueUrl, nil
}

func SendEvent(data any, queueUrl string) (string, error) {
	awsRegion := os.Getenv("AWS_REGION")

	fmt.Printf("Region: %s\n", awsRegion)
	fmt.Printf("Queue URL: %s\n", queueUrl)

	sess, err := session.NewSession(&aws.Config{
		Region:     aws.String(awsRegion),
		LogLevel:   aws.LogLevel(aws.LogDebugWithRequestErrors | aws.LogDebugWithRequestRetries | aws.LogDebugWithHTTPBody),
		MaxRetries: aws.Int(2),
		HTTPClient: &http.Client{
			Timeout: time.Duration(1 * time.Second),
		},
	},
	)

	if err != nil {
		return "", err
	}
	// fmt.Printf("Session: %s\n", *(sess.Config.Endpoint))

	svc := sqs.New(sess)
	fmt.Printf("Service ID: %s\n", svc.ServiceID)

	serial, err := serialize(data)
	fmt.Printf("Serial: %s\n", serial)

	if err != nil {
		return "", err
	}

	fmt.Println("Config Input")
	cfg := sqs.SendMessageInput{
		MessageBody: aws.String(serial),
		QueueUrl:    &queueUrl,
	}
	fmt.Println("About to send")

	result, err := svc.SendMessage(&cfg)
	fmt.Printf("MessageId: %s\n", *result.MessageId)

	return *result.MessageId, err
}

func SendEventByName(data any, queueName string) (string, error) {
	awsRegion := os.Getenv("AWS_REGION")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)

	if err != nil {
		return "", err
	}

	svc := sqs.New(sess)

	serial, err := serialize(data)

	if err != nil {
		return "", err
	}

	queueUrl, err := getQueueUrl(svc, queueName)

	if err != nil {
		return "", err
	}

	result, err := svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(serial),
		QueueUrl:    &queueUrl,
	})

	return *result.MessageId, err
}

func SendProvisionCredsEvent(data ProvisionCredsData) (string, error) {
	queueUrl := os.Getenv("PROVISION_CREDS_QUEUE_URL")

	return SendEvent(data, queueUrl)
}

func SendProvisionCredsEventByName(data ProvisionCredsData) (string, error) {
	queueName := os.Getenv("PROVISION_CREDS_QUEUE_NAME")

	return SendEventByName(data, queueName)
}

func SendSupportStaffRequestEvent(data RequestSupportStaffData) (string, error) {
	queueUrl := os.Getenv("REQUEST_SUPPORT_STAFF_QUEUE_URL")

	return SendEvent(data, queueUrl)
}

func SendSupportStaffRequestEventByName(data RequestSupportStaffData) (string, error) {
	queueName := os.Getenv("REQUEST_SUPPORT_STAFF_QUEUE_NAME")

	return SendEventByName(data, queueName)
}

func Send2FAEvent(data TwoFactorAuthData) (string, error) {
	queueUrl := os.Getenv("EMAIL_2FA_QUEUE_URL")

	return SendEvent(data, queueUrl)
}

func Send2FAEventByName(data TwoFactorAuthData) (string, error) {
	queueName := os.Getenv("EMAIL_2FA_QUEUE_NAME")

	return SendEventByName(data, queueName)
}
