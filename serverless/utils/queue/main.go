package queue

import (
	"context"

	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func SendToQueue(body string, queueUrl string) (string, error) {
	var err error
	var messageId string

	// Set up AWS configuration needed by SQS client.
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.LogError(err, "Error Loading AWS Config")
		return messageId, err
	}

	client := sqs.NewFromConfig(cfg)

	messageInput := &sqs.SendMessageInput{
		DelaySeconds: 0,
		MessageBody:  aws.String(body),
		QueueUrl:     &queueUrl,
	}

	// Send the message to SQS.
	resp, err := client.SendMessage(context.TODO(), messageInput)

	if err != nil {
		logs.LogError(err, "Failed to Send Queue Message")
		return messageId, err
	}

	messageId = *resp.MessageId

	return messageId, err
}
