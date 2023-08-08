package data

import (
	"fmt"
	"time"
)

// CheckForActiveAdmin opens a database connection and checks whether the provided
// user email exists in the `admins` table and has the `active` value set to `true`.
func CheckForActiveAdmin(adminEmail string) (bool, error) {
	var active bool
	var err error

	pool := connectToDB()
	defer pool.Close()

	query := fmt.Sprintf(`SELECT active FROM admins WHERE email = '%s';`, adminEmail)
	err = pool.QueryRow(query).Scan(&active)

	if err != nil {
		logError(err)
		return active, err
	}

	return active, err
}

// CreateAdmin opens a database connection and records the association between an admin
// user inviter and a guest user invitee along with the date of the invitation.
func CreateAdmin(adminEmail string) error {
	var err error

	pool := connectToDB()
	defer pool.Close()

	currentTime := time.Now()

	insertAdmin := `INSERT INTO "admins"("email", "active", "date_created") VALUES ($1, $2, $3);`
	_, err = pool.Exec(insertAdmin, adminEmail, true, currentTime)

	logError(err)

	return err
}
