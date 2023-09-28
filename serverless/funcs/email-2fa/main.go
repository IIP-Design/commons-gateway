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
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"

	"github.com/IIP-Design/commons-gateway/utils/logs"
)

const (
	Subject = "Verification Code for Content Commons Login"
	CharSet = "UTF-8"
)

type TwoFactorAuthData struct {
	Code  string `json:"code"`
	Email string `json:"email"`
}

// func formatEmailBody(user data.User, code string) string {
// 	return fmt.Sprintf(`<p>%s %s,</p>
// 		<p>Please use this verification code to complete your sign in:</p>
// 		<p>%s</p>
// 		<p>If you did not make this request, please disregard this email. </p>`,
// 		user.NameFirst, user.NameLast, code)
// }

// formatEmailBody constructs the body of the 2FA email.
func formatEmailBody(code string) string {
	return fmt.Sprintf(
		`<p>Please use this verification code to complete your sign in:</p>
		<p>%s</p>
		<p>If you did not make this request, please disregard this email. </p>`,
		code)
}

// formatEmail prepares the email to be sent providing a user with 2FA.
func formatEmail(email string, code string, sourceEmail string) ses.SendEmailInput {
	return ses.SendEmailInput{
		Destination: &types.Destination{
			CcAddresses: []string{},
			ToAddresses: []string{
				email,
			},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Html: &types.Content{
						Charset: aws.String(CharSet),
						Data:    aws.String(formatEmailBody(code)),
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

		emailInput := formatEmail(mfaInfo.Email, mfaInfo.Code, sourceEmail)

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
