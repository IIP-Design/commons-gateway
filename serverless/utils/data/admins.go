package data

import "time"

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
