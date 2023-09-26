package invites

import (
	"time"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/rs/xid"
)

// saveInvite opens a database connection and records the association between an admin
// user inviter and a guest user invitee along with the date of the invitation.
func SaveInvite(adminEmail string, guestEmail string, setPending bool) error {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	currentTime := time.Now()

	if setPending {
		insertInvite := `INSERT INTO invites( invitee, proposer, pending, date_invited ) VALUES ($1, $2, $3, $4);`
		_, err = pool.Exec(insertInvite, guestEmail, adminEmail, setPending, currentTime)
	} else {
		insertInvite := `INSERT INTO invites( invitee, inviter, pending, date_invited ) VALUES ($1, $2, $3, $4);`
		_, err = pool.Exec(insertInvite, guestEmail, adminEmail, setPending, currentTime)
	}

	if err != nil {
		logs.LogError(err, "Save Invite Query Error")
	}

	// Add the guest to the list of all users
	guid := xid.New()

	insertAllUsers := `INSERT INTO all_users( user_id, guest_id ) VALUES ( $1, $2 );`
	_, err = pool.Exec(insertAllUsers, guid, guestEmail)

	if err != nil {
		logs.LogError(err, "Add Guest to All Users Query Error")
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

	currentTime := time.Now()

	insertCreds :=
		`INSERT INTO guests( email, first_name, last_name, role, team, pass_hash, salt, expiration, date_created, date_modified )
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`
	_, err = pool.Exec(insertCreds, guest.Email, guest.NameFirst, guest.NameLast, guest.Role, guest.Team, hash, salt, expires, currentTime, currentTime)

	if err != nil {
		logs.LogError(err, "Save Credentials Query Error")
	}

	return err
}
