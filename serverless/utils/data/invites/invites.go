package invites

import (
	"time"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
)

// saveInvite opens a database connection and records the association between an admin
// user inviter and a guest user invitee along with the date of the invitation.
func SaveInvite(adminEmail string, guestEmail string) error {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	currentTime := time.Now()

	insertInvite := `INSERT INTO "invites"("invitee", "inviter", "date_invited") VALUES ($1, $2, $3);`
	_, err = pool.Exec(insertInvite, guestEmail, adminEmail, currentTime)

	if err != nil {
		logs.LogError(err, "Save Invite Query Error")
	}

	return err
}

// SaveCredentials opens a database connection and saves the provided user credentials
// to the `credentials` table. Specifically, it stores the the user email, a hash of
// their password, and the salt with which the password was hashed, as well as the date
// on which the password was generated.
func SaveCredentials(guest data.User, hash string, salt string) error {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	currentTime := time.Now()
	// TODO: Placeholder value; expiration this should be set via the client app
	expiration := currentTime.AddDate(0, 2, 0)

	insertCreds :=
		`INSERT INTO "guests"("email", "first_name", "last_name", "team", "pass_hash", "salt", "expiration", "date_created" ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`
	_, err = pool.Exec(insertCreds, guest.Email, guest.NameFirst, guest.NameLast, guest.Team, hash, salt, expiration, currentTime)

	if err != nil {
		logs.LogError(err, "Save Credentials Query Error")
	}

	return err
}
