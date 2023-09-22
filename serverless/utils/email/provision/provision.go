package provision

import (
	"fmt"
	"log"
	"os"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	email "github.com/IIP-Design/commons-gateway/utils/email/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	Subject = "Content Commons Account Created"
	CharSet = "UTF-8"
)

type ProvisionCredsData struct {
	Invitee     data.User `json:"invitee"`
	TmpPassword string    `json:"tmpPassword"`
	Url         string    `json:"url"`
}

func formatEmailBody(invitee data.User, tmpPassword string, url string) string {
	return fmt.Sprintf(`<p>%s %s,</p>

	<p>Your content upload account has been successfully created.  Please access the link below to finish provisioning your account.</p>
	<a href="%s">%s</a>
	<p>Please use this email address as your username.  Your temporary password is: %s.</p>
	<p>This email was generated automatically. Please do not reply to this email.</p>`,
		invitee.NameFirst, invitee.NameLast,
		url, url,
		tmpPassword)
}

func formatEmail(invitee data.User, tmpPassword string, url string, sourceEmail string) ses.SendEmailInput {
	return ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(invitee.Email),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(formatEmailBody(invitee, tmpPassword, url)),
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

func MailProvisionedCreds(sourceEmail string, provisionCredsData ProvisionCredsData) error {
	if sourceEmail == "" {
		log.Println("Not configured for sending emails")
		return nil
	}

	region := os.Getenv("AWS_SES_REGION")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return err
	}

	sesClient := ses.New(sess)

	e := formatEmail(
		provisionCredsData.Invitee,
		provisionCredsData.TmpPassword,
		provisionCredsData.Url,
		sourceEmail,
	)

	result, err := sesClient.SendEmail(&e)
	email.LogSesResult(result, err)

	return nil
}
