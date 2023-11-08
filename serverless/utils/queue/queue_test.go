package queue

import (
	"context"
	"testing"

	testConfig "github.com/IIP-Design/commons-gateway/test/config"
	testHelpers "github.com/IIP-Design/commons-gateway/test/helpers"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

const (
	QUEUE_NAME   = "test_queue"
	MESSAGE_BODY = "test"
)

func TestSendToQueue(t *testing.T) {
	testConfig.ConfigureAws()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		t.Fatalf("AWS config error: %v", err)
	}

	client := sqs.NewFromConfig(cfg)

	_, err = testHelpers.CreateQueue(QUEUE_NAME, client)
	if err != nil {
		t.Fatalf("CreateQueue error: %v", err)
	}

	queueUrl, err := testHelpers.GetQueueUrl(QUEUE_NAME, client)
	if err != nil {
		t.Fatalf("GetQueueUrl error: %v", err)
	}

	sendId, err := SendToQueue(MESSAGE_BODY, queueUrl)
	if err != nil {
		t.Fatalf("SendToQueue error: %v", err)
	}

	msg, err := testHelpers.GetMessages(QUEUE_NAME, client)
	if err != nil {
		t.Fatalf("GetMessages error: %v", err)
	}

	msgBody := msg.Messages[0].Body
	msgId := msg.Messages[0].MessageId

	if *msgBody != MESSAGE_BODY {
		t.Fatalf("Message body error: %s, want %s", *msgBody, MESSAGE_BODY)
	}
	if *msgId != sendId {
		t.Fatalf("Message ID error: %s, want %s", *msgId, sendId)
	}

	err = testHelpers.DeleteQueue(QUEUE_NAME, client)
	if err != nil {
		t.Fatalf("DeleteQueue error: %v", err)
	}
}
