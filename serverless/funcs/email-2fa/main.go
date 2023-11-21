package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	ses "github.com/aws/aws-sdk-go-v2/service/sesv2"
	sesTypes "github.com/aws/aws-sdk-go-v2/service/sesv2/types"

	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/IIP-Design/commons-gateway/utils/types"
)

const (
	Subject = "Verification Code for Content Commons Login"
	CharSet = "UTF-8"
)

type TwoFactorAuthData struct {
	Code string     `json:"code"`
	User types.User `json:"user"`
}

// formatEmailBody constructs the body of the 2FA email.
func formatEmailBody(user types.User, code string) string {
	return fmt.Sprintf(
		`<p>%s %s,</p>
		<p>Please use this verification code to complete your sign in:</p>
		<p>%s</p>
		<p>Please note that this verification code will expire in 20 minutes. If you did not make this request, please disregard this email. </p>`,
		user.NameFirst, user.NameLast, code,
	)
}

// formatEmail prepares the email to be sent providing a user with 2FA.
func formatEmail(user types.User, code string, sourceEmail string) ses.SendEmailInput {
	return ses.SendEmailInput{
		Destination: &sesTypes.Destination{
			CcAddresses: []string{},
			ToAddresses: []string{
				user.Email,
			},
		},
		Content: &sesTypes.EmailContent{
			Simple: &sesTypes.Message{
				Body: &sesTypes.Body{
					Html: &sesTypes.Content{
						Charset: aws.String(CharSet),
						Data:    aws.String(formatEmailBody(user, code)),
					},
				},
				Subject: &sesTypes.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(Subject),
				},
			},
		},
		FromEmailAddress: &sourceEmail,
	}
}

// email2FAHandler sends a guest user a 2FA code that they can use to log in.
func email2FAHandler(ctx context.Context, event events.SQSEvent) error {
	var err error

	region := os.Getenv("AWS_SES_REGION")
	sourceEmail := os.Getenv("SOURCE_EMAIL_ADDRESS")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)

	if err != nil {
		return err
	}

	sesClient := ses.NewFromConfig(cfg)

	for _, message := range event.Records {
		eventMessageId := message.MessageId
		body := message.Body

		var mfaInfo TwoFactorAuthData

		err := json.Unmarshal([]byte(body), &mfaInfo)

		if err != nil {
			logs.LogError(err, fmt.Sprintf("Unable to unmarshal body of message %s", eventMessageId))
			return err
		}

		emailInput := formatEmail(mfaInfo.User, mfaInfo.Code, sourceEmail)

		_, err = sesClient.SendEmail(context.TODO(), &emailInput)
		if err != nil {
			log.Println(err.Error())
		}
	}

	return nil
}

func main() {
	lambda.Start(email2FAHandler)
}
