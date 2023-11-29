package propose

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
	Subject = "Content Commons Support Staff Request"
	CharSet = "UTF-8"
)

type RequestSupportStaffData struct {
	Proposer data.User `json:"externalTeamLead"`
	Invitee  data.User `json:"supportStaffuser"`
	Url      string    `json:"url"`
}

// getAdmins retrieves a list of admins assigned to a given team.
func getAdmins(team string) ([]data.User, error) {
	var admins []data.User

	pool := data.ConnectToDB()
	defer pool.Close()

	query := "SELECT email, first_name, last_name, role, team FROM admins WHERE team = $1"
	rows, err := pool.Query(query, team)

	if err != nil {
		logs.LogError(err, "Get Uploaders Query Error")
		return admins, err
	}

	defer rows.Close()

	for rows.Next() {
		var admin data.User
		err := rows.Scan(
			&admin.Email,
			&admin.NameFirst,
			&admin.NameLast,
			&admin.Role,
			&admin.Team,
		)

		if err != nil {
			logs.LogError(err, "Get Admins Scan Error")
			return admins, err
		}

		admins = append(admins, admin)
	}

	if err = rows.Err(); err != nil {
		logs.LogError(err, "Get Admins Row Error")
		return admins, err
	}

	return admins, nil
}

// formatEmailBody populates the template of the invite proposal
// notification email with the relevant values.
func formatEmailBody(
	proposer data.User,
	invitee data.User,
	admin data.User,
	url string,
) string {
	return fmt.Sprintf(
		`<p>%s %s,</p>
		 <p>%s %s has submitted a ticket for adding %s %s for your approval. Please follow <a href="%s">this link</a> to approve or deny this request.</p>
		 <p>This email was generated automatically. Please do not reply to this email.</p>`,
		admin.NameFirst, admin.NameLast,
		proposer.NameFirst, proposer.NameLast,
		invitee.NameFirst, invitee.NameLast,
		url,
	)
}

// formatEmail prepares the event object that is sent to SES in order
// to initiate the invite proposal notification email.
func formatEmail(
	proposer data.User,
	invitee data.User,
	admin data.User,
	url string,
	sourceEmail string,
) ses.SendEmailInput {
	return ses.SendEmailInput{
		Destination: &sesTypes.Destination{
			CcAddresses: []string{},
			ToAddresses: []string{
				admin.Email,
			},
		},
		Content: &sesTypes.EmailContent{
			Simple: &sesTypes.Message{
				Body: &sesTypes.Body{
					Html: &sesTypes.Content{
						Charset: aws.String(CharSet),
						Data:    aws.String(formatEmailBody(proposer, invitee, admin, url)),
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

func MailProposedCreds(proposer data.User, invitee data.User) error {
	sourceEmail := os.Getenv("SOURCE_EMAIL_ADDRESS")
	redirectUrl := os.Getenv("EMAIL_REDIRECT_URL")

	if sourceEmail == "" {
		err := errors.New("not configured for sending emails")

		logs.LogError(err, "Missing Source Email Error")
		return nil
	}

	awsRegion := os.Getenv("AWS_SES_REGION")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
	)

	if err != nil {
		logs.LogError(err, "Error Loading AWS Configuration")
		return err
	}

	admins, err := getAdmins(proposer.Team)

	if err != nil {
		logs.LogError(err, "Get Admins Error")
		return err
	}

	sesClient := ses.NewFromConfig(cfg)

	for _, admin := range admins {
		e := formatEmail(
			proposer,
			invitee,
			admin,
			redirectUrl,
			sourceEmail,
		)

		_, err := sesClient.SendEmail(context.TODO(), &e)

		if err != nil {
			logs.LogError(err, "Send Error")
		}

	}

	return nil
}
