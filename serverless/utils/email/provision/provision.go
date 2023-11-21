package provision

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	ses "github.com/aws/aws-sdk-go-v2/service/sesv2"
	sesTypes "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

const (
	CharSet = "UTF-8"
)

// There are various cases in which we provision credentials.
// This type enumerates those situations.
type ProvisionType int

const (
	Create ProvisionType = iota
	Reauth
	Reset
)

// Subject returns the email subject line appropriate to
// a given provisioning action.
func (pt ProvisionType) Subject() string {
	switch pt {
	case Create:
		return "Content Commons Account Created"
	case Reauth:
		return "Content Commons Account Reactivation"
	case Reset:
		return "Content Commons Password Reset"
	default:
		return "Content Commons Account"
	}
}

// Verb returns the email subject line appropriate to
// a given provisioning action.
func (pt ProvisionType) Verb() string {
	switch pt {
	case Create:
		return "created"
	case Reauth:
		return "reactivated"
	case Reset:
		return "reset"
	default:
		return "updated"
	}
}

func formatEmailBody(invitee data.User, tmpPassword string, url string, verb string) string {
	return fmt.Sprintf(
		`<p>%s %s,</p>
		<p>Your content upload account has been successfully %s. Please access the link below to finish provisioning your account.</p>
		<a href="%s">%s</a>
		<p>Please use this email address as your username. Your temporary password is: %s</p>
		<p>This email was generated automatically. Please do not reply to this email.</p>`,
		invitee.NameFirst,
		invitee.NameLast,
		verb,
		url,
		url,
		tmpPassword,
	)
}

func formatEmail(invitee data.User, tmpPassword string, url string, sourceEmail string, action ProvisionType) ses.SendEmailInput {

	return ses.SendEmailInput{
		Destination: &sesTypes.Destination{
			CcAddresses: []string{},
			ToAddresses: []string{
				invitee.Email,
			},
		},
		Content: &sesTypes.EmailContent{
			Simple: &sesTypes.Message{
				Body: &sesTypes.Body{
					Html: &sesTypes.Content{
						Charset: aws.String(CharSet),
						Data:    aws.String(formatEmailBody(invitee, tmpPassword, url, action.Verb())),
					},
				},
				Subject: &sesTypes.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(action.Subject()),
				},
			},
		},
		FromEmailAddress: &sourceEmail,
	}
}

// MailProvisionedCreds emails the user a temporary password that can be used to login
// into the external partner portal. For the action parameter, pass in an integer corresponding
// to one of the credential provisioning actions. There are three enumerated action types:
//
//	1 - used when creating a new account
//	2 - used when reauthorizing an existing expired account
//	3 - used when resetting an existing account password
func MailProvisionedCreds(invitee data.User, tmpPassword string, action ProvisionType) (string, error) {
	var err error
	var mesageId string

	sourceEmail := os.Getenv("SOURCE_EMAIL_ADDRESS")
	redirectUrl := os.Getenv("EMAIL_REDIRECT_URL")

	if sourceEmail == "" {
		logs.LogError(errors.New("not configured for sending emails"), "Source Email Empty Error")
		return mesageId, err
	}

	awsRegion := os.Getenv("AWS_SES_REGION")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))

	if err != nil {
		logs.LogError(err, "Error Loading AWS Config")
		return mesageId, err
	}

	sesClient := ses.NewFromConfig(cfg)

	e := formatEmail(
		invitee,
		tmpPassword,
		redirectUrl,
		sourceEmail,
		action,
	)

	resp, err := sesClient.SendEmail(context.TODO(), &e)
	if err != nil {
		logs.LogError(err, "Credentials Provisioning Email Error")
	}

	mesageId = *resp.MessageId

	return mesageId, err
}
