package data

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/IIP-Design/commons-gateway/utils/logs"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	_ "github.com/lib/pq"
)

const (
	aws_region  = "DB_REGION"
	db_host     = "DB_HOST"
	db_name     = "DB_NAME"
	db_password = "DB_PASSWORD"
	db_port     = "DB_PORT"
	db_user     = "DB_USER"
)

// ConnectToDB opens a pool connection to a Postgres database using
// credentials derived from the environment.
func ConnectToDB() *sql.DB {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		logs.LogError(err, "DB Configuration Error")
	}

	host := os.Getenv(db_host)
	name := os.Getenv(db_name)
	password := os.Getenv(db_password)
	port := os.Getenv(db_port)
	user := os.Getenv(db_user)
	region := os.Getenv(aws_region)

	var connStr string

	// IAM Authentication is required when the app is deployed to AWS.
	// A password-based authentication option is preserved for local function testing.
	if password == "" {
		authToken, err := auth.BuildAuthToken(
			context.TODO(),
			fmt.Sprintf("%s:%s", host, port),
			region,
			user,
			cfg.Credentials,
		)

		if err != nil {
			logs.LogError(err, "DB Authentication Token Error")
		}

		connStr = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s",
			host,
			port,
			user,
			authToken,
			name,
		)
	} else {
		// Connection string for local testing.
		connStr = fmt.Sprintf(
			"postgresql://%s:%s@%s/%s?sslmode=disable",
			user,
			password,
			host,
			name,
		)
	}

	pool, err := sql.Open("postgres", connStr)

	if err != nil {
		logs.LogError(err, "DB Connection Error")
	}

	return pool
}
