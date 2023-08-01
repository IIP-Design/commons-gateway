package data

import (
	"database/sql"
	"fmt"
	"os"
	"time"

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

// saveInvite opens a database connection and records the association between an admin
// user inviter and a guest user invitee along with the date of the invitation.
func SaveInvite(adminEmail string, guestEmail string) error {
	var err error

	pool := connectToDB()

	defer pool.Close()

	currentTime := time.Now()

	insertInvite := `INSERT INTO "invites"("invitee", "inviter", "date_invited") VALUES ($1, $2, $3);`
	_, err = pool.Exec(insertInvite, adminEmail, guestEmail, currentTime)

	logError(err)

	return err
}

// saveCredentials opens a database connection and saves the provided user credentials
// to the `credentials` table. Specifically, it stores the the user email, a hash of
// their password, and the salt with which the password was hashed, as well as the date
// on which the password was generated.
func SaveCredentials(email string, hash string, salt string) error {
	var err error

	pool := connectToDB()

	defer pool.Close()

	currentTime := time.Now()

	insertCreds := `INSERT INTO "credentials"("email", "pass_hash", "salt", "date_created" ) VALUES ($1, $2, $3, $4);`
	_, err = pool.Exec(insertCreds, email, hash, salt, currentTime)

	logError(err)

	return err
}

// checkForExistingUser opens a database connection and checks whether the provided
// email (which is a unique value constraint in the admins and credentials tables) is
// present in the provided table. An affirmative check indicates that the given user
// has the access implied by their presence in the table.
func CheckForExistingUser(email string, table string) (bool, error) {
	var err error
	var hasAccess = false

	pool := connectToDB()

	query := fmt.Sprintf(`SELECT "email" FROM %s;`, table)

	rows, err := pool.Query(query)
	logError(err)

	for rows.Next() {
		var retrievedEmail string

		err = rows.Scan(&retrievedEmail)
		logError(err)

		if retrievedEmail == email {
			hasAccess = true
			break
		}
	}

	return hasAccess, err
}
