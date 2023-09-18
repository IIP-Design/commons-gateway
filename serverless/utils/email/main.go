package email

import (
	"bytes"
	"encoding/json"
	"os"

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
	Inviter  data.User `json:"contentCommonsUser"`
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

func send(data any, queueUrl string) (string, error) {
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

	result, err := svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(serial),
		QueueUrl:    &queueUrl,
	})

	return *result.MessageId, err
}

func SendProvisionCredsEvent(data ProvisionCredsData) (string, error) {
	queueUrl := os.Getenv("PROVISION_CREDS_QUEUE")

	return send(data, queueUrl)
}

func SendSupportStaffRequestEvent(data RequestSupportStaffData) (string, error) {
	queueUrl := os.Getenv("REQUEST_SUPPORT_STAFF_QUEUE")

	return send(data, queueUrl)
}

func Send2FAEvent(data TwoFactorAuthData) (string, error) {
	queueUrl := os.Getenv("EMAIL_2FA_QUEUE")

	return send(data, queueUrl)
}
