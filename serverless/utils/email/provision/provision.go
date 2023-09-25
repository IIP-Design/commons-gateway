package provision

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	ses "github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
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
		Destination: &types.Destination{
			CcAddresses: []string{},
			ToAddresses: []string{
				invitee.Email,
			},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Html: &types.Content{
						Charset: aws.String(CharSet),
						Data:    aws.String(formatEmailBody(invitee, tmpPassword, url)),
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

func MailProvisionedCreds(sourceEmail string, provisionCredsData ProvisionCredsData) error {
	if sourceEmail == "" {
		log.Println("Not configured for sending emails")
		return nil
	}

	awsRegion := os.Getenv("AWS_SES_REGION")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
	)
	if err != nil {
		return err
	}

	sesClient := ses.NewFromConfig(cfg)

	e := formatEmail(
		provisionCredsData.Invitee,
		provisionCredsData.TmpPassword,
		provisionCredsData.Url,
		sourceEmail,
	)

	_, err = sesClient.SendEmail(context.TODO(), &e)
	if err != nil {
		log.Println(err.Error())
	}

	return nil
}
