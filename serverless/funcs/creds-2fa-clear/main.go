package main

import (
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"

	"github.com/aws/aws-lambda-go/lambda"
)

// expiredMFACodeHandler deletes any 2FA request older than 20 minutes.
func expiredMFACodeHandler() {
	pool := data.ConnectToDB()
	defer pool.Close()

	query := `DELETE FROM mfa WHERE date_created < now()::timestamp - INTERVAL '20 minutes';`
	_, err := pool.Exec(query)

	if err != nil {
		logs.LogError(err, "Clear Expired MFA Code Query Error")
	}
}

func main() {
	lambda.Start(expiredMFACodeHandler)
}
