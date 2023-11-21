package testHelpers

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func GetQueueUrl(queueName string, client *sqs.Client) (string, error) {
	result, err := client.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})

	if err == nil {
		return *result.QueueUrl, err
	} else {
		return "", err
	}
}

func CreateQueue(queueName string, client *sqs.Client) (*sqs.CreateQueueOutput, error) {
	return client.CreateQueue(context.TODO(), &sqs.CreateQueueInput{
		QueueName: &queueName,
		Attributes: map[string]string{
			"MessageRetentionPeriod": "86400",
		},
	})
}

func DeleteQueue(queueName string, client *sqs.Client) error {
	queueUrl, err := GetQueueUrl(queueName, client)
	if err != nil {
		return err
	}

	_, err = client.DeleteQueue(context.TODO(), &sqs.DeleteQueueInput{
		QueueUrl: &queueUrl,
	})

	return err
}

func GetMessages(queueName string, client *sqs.Client) (*sqs.ReceiveMessageOutput, error) {
	queueUrl, err := GetQueueUrl(queueName, client)
	if err != nil {
		return nil, err
	}

	return client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:            &queueUrl,
		MaxNumberOfMessages: 1,
	})
}

func CreateTestQueue(queueName string, envName string) (string, *sqs.Client, error) {
	var queueUrl string
	var client *sqs.Client
	var err error

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return queueUrl, client, err
	}

	client = sqs.NewFromConfig(cfg)

	_, err = CreateQueue(queueName, client)
	if err != nil {
		return queueUrl, client, err
	}

	queueUrl, err = GetQueueUrl(queueName, client)
	if err != nil {
		return queueUrl, client, err
	}

	if envName != "" {
		_ = os.Setenv(envName, queueUrl)
	}

	return queueUrl, client, err
}
