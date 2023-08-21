package data

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/IIP-Design/commons-gateway/utils/logs"
	_ "github.com/lib/pq"
)

const (
	db_host     = "DB_HOST"
	db_name     = "DB_NAME"
	db_password = "DB_PASSWORD"
	db_user     = "DB_USER"
)

// ConnectToDB opens a pool connection to a Postgres database using
// credentials derived from the environment.
func ConnectToDB() *sql.DB {
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

	if err != nil {
		logs.LogError(err, "DB Connection Error")
	}

	return pool
}

// CheckForExistingUser opens a database connection and checks whether the provided
// email (which is a unique value constraint in the admins and guests tables) is
// present in the provided table. An affirmative check indicates that the given user
// has the access implied by their presence in the table.
func CheckForExistingUser(email string, table string) (bool, error) {
	var err error

	pool := ConnectToDB()
	defer pool.Close()

	var user string

	query := fmt.Sprintf(`SELECT email FROM %s WHERE email = '%s';`, table, email)
	err = pool.QueryRow(query).Scan(&user)

	if err != nil {
		// Do not return an error if no results are found.
		if err == sql.ErrNoRows {
			return false, nil
		}

		logs.LogError(err, "Existing User Query Error")
	}

	return user == email, err
}
