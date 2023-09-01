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

	// For now all invites will be set to not pending,
	// will change when invite proposals are introduced
	pending := false
	currentTime := time.Now()

	insertInvite := `INSERT INTO invites( invitee, inviter, pending, date_invited ) VALUES ($1, $2, $3, $4);`
	_, err = pool.Exec(insertInvite, guestEmail, adminEmail, pending, currentTime)

	if err != nil {
		logs.LogError(err, "Save Invite Query Error")
	}

	return err
}

// SaveCredentials opens a database connection and saves the provided user credentials
// to the `credentials` table. Specifically, it stores the the user email, a hash of
// their password, and the salt with which the password was hashed, as well as the date
// on which the password was generated.
func SaveCredentials(guest data.User, expires time.Time, hash string, salt string) error {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	role := "guest"
	currentTime := time.Now()

	insertCreds :=
		`INSERT INTO guests( email, first_name, last_name, role, team, pass_hash, salt, expiration, date_created, date_modified )
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`
	_, err = pool.Exec(insertCreds, guest.Email, guest.NameFirst, guest.NameLast, role, guest.Team, hash, salt, expires, currentTime, currentTime)

	if err != nil {
		logs.LogError(err, "Save Credentials Query Error")
	}

	return err
}
