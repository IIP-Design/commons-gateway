package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	ses "github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

const (
	Subject = "Verification Code for Content Commons Login"
	CharSet = "UTF-8"
)

type TwoFactorAuthData struct {
	User data.User `json:"user"`
	Code string    `json:"verificationCode"`
}

func formatEmailBody(
	user data.User,
	code string,
) string {
	return fmt.Sprintf(`<p>%s %s,</p>
		<p>Please use this verification code to complete your sign in:</p>
		<p>%s</p>
		<p>If you did not make this request, please disregard this email. </p>`,
		user.NameFirst, user.NameLast, code)
}

func formatEmail(
	user data.User,
	code string,
	sourceEmail string,
) ses.SendEmailInput {
	return ses.SendEmailInput{
		Destination: &types.Destination{
			CcAddresses: []string{},
			ToAddresses: []string{
				user.Email,
			},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Html: &types.Content{
						Charset: aws.String(CharSet),
						Data:    aws.String(formatEmailBody(user, code)),
					},
				},
				Subject: &types.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(Subject),
				},
			},
		},
		FromEmailAddress: &sourceEmail,
	}
}

func Email2FAHandler(ctx context.Context, event events.SQSEvent) error {
	var err error
	records := event.Records

	region := os.Getenv("AWS_SES_REGION")
	sourceEmail := os.Getenv("SOURCE_EMAIL_ADDRESS")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return err
	}

	sesClient := ses.NewFromConfig(cfg)

	for _, record := range records {
		eventMessageId := record.MessageId
		body := record.Body

		var userData TwoFactorAuthData
		err = json.Unmarshal([]byte(body), &userData)

		if err != nil {
			logs.LogError(err, fmt.Sprintf("Unable to unmarshal body of message %s", eventMessageId))
			return err
		}

		e := formatEmail(userData.User, userData.Code, sourceEmail)

		_, err = sesClient.SendEmail(context.TODO(), &e)
		if err != nil {
			log.Println(err.Error())
		}
	}

	return nil
}

func main() {
	lambda.Start(Email2FAHandler)
}
