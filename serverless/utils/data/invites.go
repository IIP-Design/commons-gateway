package data

import "time"

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

// SaveCredentials opens a database connection and saves the provided user credentials
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
