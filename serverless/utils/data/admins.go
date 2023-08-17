package data

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/IIP-Design/commons-gateway/utils/logs"
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
		logs.LogError(err, "Existing Admin Query Error")
		return active, err
	}

	return active, err
}

// CreateAdmin opens a database connection and records the association between an admin
// user inviter and a guest user invitee along with the date of the invitation.
func CreateAdmin(adminData User) error {
	var err error

	pool := connectToDB()
	defer pool.Close()

	currentTime := time.Now()

	insertAdmin :=
		`INSERT INTO "admins"("email", "first_name", "last_name", "team", "active", "date_created")
		 VALUES ($1, $2, $3, $4, $5, $6);`
	_, err = pool.Exec(insertAdmin, adminData.Email, adminData.NameFirst, adminData.NameLast, adminData.Team, true, currentTime)

	if err != nil {
		logs.LogError(err, "Create Admin Query Error")
	}

	return err
}
