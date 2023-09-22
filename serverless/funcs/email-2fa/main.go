package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	email "github.com/IIP-Design/commons-gateway/utils/email/utils"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
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
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(user.Email),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(formatEmailBody(user, code)),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(sourceEmail),
	}
}

func Email2FAHandler(ctx context.Context, event events.SQSEvent) error {
	var err error
	records := event.Records

	region := os.Getenv("AWS_SES_REGION")
	sourceEmail := os.Getenv("SOURCE_EMAIL_ADDRESS")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return err
	}

	svc := ses.New(sess)

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

		result, err := svc.SendEmail(&e)
		email.LogSesResult(result, err)
	}

	return nil
}

func main() {
	lambda.Start(Email2FAHandler)
}
