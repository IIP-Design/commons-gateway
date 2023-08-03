package data

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const (
	db_host     = "DB_HOST"
	db_name     = "DB_NAME"
	db_password = "DB_PASSWORD"
	db_user     = "DB_USER"
)

func logError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// connectToDB opens a pool connection to a Postgres database using
// credentials derived from the environment.
func connectToDB() *sql.DB {
	host := os.Getenv(db_host)
	name := os.Getenv(db_name)
	password := os.Getenv(db_password)
	user := os.Getenv(db_user)

	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s",
		user,
		password,
		host,
		name,
	)

	pool, err := sql.Open("postgres", connStr)
	logError(err)

	return pool
}

// checkForExistingUser opens a database connection and checks whether the provided
// email (which is a unique value constraint in the admins and credentials tables) is
// present in the provided table. An affirmative check indicates that the given user
// has the access implied by their presence in the table.
func CheckForExistingUser(email string, table string) (bool, error) {
	var exists bool
	var err error

	pool := connectToDB()

	var found string

	query := fmt.Sprintf(`SELECT 1 FROM %s WHERE email = '%s';`, table, email)
	err = pool.QueryRow(query).Scan(&found)

	if err != nil {
		logError(err)
		return exists, err
	}

	return found == "1", err
}
