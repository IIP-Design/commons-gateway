package queue

import (
	"context"
	"os"

	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func SendToQueue(body string, queueUrl string) (string, error) {
	var cfg aws.Config
	var err error
	var messageId string

	dbg := os.Getenv("DEBUG") == "true"

	// Set up AWS configuration needed by SQS client.
	if dbg {
		cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithClientLogMode(aws.LogRetries|aws.LogRequestWithBody|aws.LogResponseWithBody|aws.LogRequestEventMessage))
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO())
	}

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
