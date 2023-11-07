package init

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/rs/xid"
)

// Copied from data module to prevent circular import(s) during testing
func connectToDB() *sql.DB {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		logs.LogError(err, "DB Configuration Error")
	}

	host := os.Getenv("DB_HOST")
	name := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	region := os.Getenv("DB_REGION")

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

// CheckForTable identifies whether or not a table with a given name
// is present in the database.
func CheckForTable(tablename string) bool {
	var exists bool

	pool := connectToDB()
	defer pool.Close()

	query :=
		`SELECT EXISTS ( SELECT FROM pg_tables WHERE schemaname = 'public'
		 AND tablename  = $1 );`

	err := pool.QueryRow(query, tablename).Scan(&exists)

	if err == sql.ErrNoRows {
		msg := fmt.Sprintf("Table Not Found - %s", tablename)

		logs.LogError(err, msg)
		return false
	} else if err != nil {
		msg := fmt.Sprintf("Check Table Query Error - %s", tablename)

		logs.LogError(err, msg)
		return false
	}

	return exists
}

// stringArrayContains iterates through an array of string to check whether
// a given string value is present as an item in the array.
func stringArrayContains(arr []string, value string) bool {
	for _, item := range arr {
		if item == value {
			return true
		}
	}

	return false
}

// recordMigration saves the title of an schema migration and the date on which it was applied.
func recordMigration(title string) error {
	pool := connectToDB()
	defer pool.Close()

	guid := xid.New()
	currentTime := time.Now()

	query := `INSERT INTO migrations( id, title, date_applied ) VALUES ( $1, $2, $3 );`

	_, err := pool.Exec(query, guid, title, currentTime)

	if err != nil {
		logs.LogError(err, "Migration Registry Error")
	}

	return err
}
