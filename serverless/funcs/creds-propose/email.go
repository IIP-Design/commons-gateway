package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	ses "github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
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

func getAdmins(team string) ([]data.User, error) {
	var admins []data.User

	pool := data.ConnectToDB()
	defer pool.Close()

	query := "SELECT email, first_name, last_name, role, team FROM admins WHERE team = $1 OR role='super admin'"
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

func formatEmailBody(
	proposer data.User,
	invitee data.User,
	admin data.User,
	url string,
) string {
	return fmt.Sprintf(`<p>%s %s,</p> 
	<p>%s %s has submitted a ticket for adding
	 %s %s for your approval.
	  Please follow <a href="%s">this link</a> to approve or deny this request.</p>
	<p>This email was generated automatically. Please do not reply to this email.</p>`,
		admin.NameFirst, admin.NameLast,
		proposer.NameFirst, proposer.NameLast,
		invitee.NameFirst, invitee.NameLast,
		url)
}

func formatEmail(
	proposer data.User,
	invitee data.User,
	admin data.User,
	url string,
	sourceEmail string,
) ses.SendEmailInput {
	return ses.SendEmailInput{
		Destination: &types.Destination{
			CcAddresses: []string{},
			ToAddresses: []string{
				admin.Email,
			},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Html: &types.Content{
						Charset: aws.String(CharSet),
						Data:    aws.String(formatEmailBody(proposer, invitee, admin, url)),
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

func MailProposedCreds(sourceEmail string, supportStaffRequestData RequestSupportStaffData) error {
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

	admins, err := getAdmins(supportStaffRequestData.Proposer.Team)
	if err != nil {
		return err
	}

	sesClient := ses.NewFromConfig(cfg)

	for _, admin := range admins {
		e := formatEmail(
			supportStaffRequestData.Proposer,
			supportStaffRequestData.Invitee,
			admin,
			supportStaffRequestData.Url,
			sourceEmail,
		)

		_, err := sesClient.SendEmail(context.TODO(), &e)
		if err != nil {
			log.Println(err.Error())
		}

	}

	return nil
}
