package guests

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/invites"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/IIP-Design/commons-gateway/utils/security/hashing"
)

type InviteRecord struct {
	Proposer      string `json:"proposer"`
	Pending       bool   `json:"pending"`
	DateInvited   string `json:"dateInvited"`
	Expiration    string `json:"expiration"`
	Expired       bool   `json:"expired"`
	Inviter       string `json:"inviter"`
	PasswordReset bool   `json:"passwordReset"`
}

type GuestData struct {
	Email     string `json:"email"`
	FirstName string `json:"givenName"`
	LastName  string `json:"familyName"`
	Role      string `json:"role"`
	Team      string `json:"team"`
}

type GuestDetails struct {
	GuestData
	Invites []InviteRecord `json:"invites"`
}

// RetrieveGuest opens a database connection and retrieves the information for a single user.
func RetrieveGuest(email string) (GuestDetails, error) {
	var guest GuestDetails

	pool := data.ConnectToDB()
	defer pool.Close()

	query := `SELECT email, first_name, last_name, role, team FROM guests WHERE email = $1`
	err := pool.QueryRow(query, email).Scan(&guest.Email, &guest.FirstName, &guest.LastName, &guest.Role, &guest.Team)

	if err != nil {
		logs.LogError(err, "Retrieve Guest Query Error")
		return guest, err
	}

	query = `SELECT i.pending, i.date_invited, i.expiration, i.expiration < NOW() AS expired, i.password_reset,
	  a.first_name, a.last_name, COALESCE( i.proposer, '' )
		FROM invites i
		JOIN admins a ON i.inviter = a.email
		WHERE invitee = $1 ORDER BY date_invited DESC`
	rows, err := pool.Query(query, email)

	if err != nil {
		logs.LogError(err, "Retrieve Invites Query Error")
		return guest, err
	}

	defer rows.Close()

	for rows.Next() {
		var inviterFirst string
		var inviterLast string
		var proposer string
		var pending bool
		var dateInvited time.Time
		var expiration time.Time
		var expired bool
		var passwordReset bool

		if err := rows.Scan(
			&pending,
			&dateInvited,
			&expiration,
			&expired,
			&passwordReset,
			&inviterFirst,
			&inviterLast,
			&proposer,
		); err != nil {
			logs.LogError(err, "Scan Guests Query Error")
			return guest, err
		}

		proposerName := ""

		if proposer != "" {
			query := "SELECT first_name, last_name FROM guests WHERE email = $1"

			var proposerFirstName string
			var proposerLastName string

			pool.QueryRow(query, proposer).Scan(&proposerFirstName, &proposerLastName)

			if err != nil || proposerFirstName == "" || proposerLastName == "" {
				logs.LogError(err, "Retrieve Proposer Query Error")
			} else {
				proposerName = fmt.Sprintf("%s %s", proposerFirstName, proposerLastName)
			}
		}

		var invite = InviteRecord{
			DateInvited:   dateInvited.Format(time.RFC3339),
			Expiration:    expiration.Format(time.RFC3339),
			Expired:       expired,
			Inviter:       fmt.Sprintf("%s %s", inviterFirst, inviterLast),
			PasswordReset: passwordReset,
			Pending:       pending,
			Proposer:      proposerName,
		}

		guest.Invites = append(guest.Invites, invite)
	}

	return guest, err
}

// RetrieveGuestExpiration opens a database connection and retrieves a single user's access expiration.
func RetrieveGuestExpiration(email string) (time.Time, error) {
	var expires time.Time

	pool := data.ConnectToDB()
	defer pool.Close()

	query := `SELECT expiration FROM invites WHERE invitee = $1 ORDER BY date_invited DESC LIMIT 1`
	err := pool.QueryRow(query, email).Scan(&expires)

	if err != nil {
		logs.LogError(err, "Retrieve Guest Expiration Query Error")
	}

	return expires, err
}

// RetrieveGuests opens a database connection and retrieves the full list of guest users.
func RetrieveGuests(team string, role string) ([]data.GuestUser, error) {
	var guests []data.GuestUser
	var err error
	var query string
	var rows *sql.Rows

	pool := data.ConnectToDB()
	defer pool.Close()

	if team == "" {
		query = `SELECT email, first_name, last_name, role, team, pending, expiration FROM guest_auth_data ORDER BY first_name;`
		rows, err = pool.Query(query)
	} else {
		query =
			`SELECT email, first_name, last_name, role, team, pending, expiration
			 FROM guest_auth_data WHERE team = $1 ORDER BY first_name;`
		rows, err = pool.Query(query, team)
	}

	if err != nil {
		logs.LogError(err, "Get Guests Query Error")
		return guests, err
	}

	defer rows.Close()

	for rows.Next() {
		var guest data.GuestUser
		if err := rows.Scan(&guest.Email, &guest.NameFirst, &guest.NameLast, &guest.Role, &guest.Team, &guest.Pending, &guest.Expires); err != nil {
			logs.LogError(err, "Get Guests Query Error")
			return guests, err
		}

		if role == "" || role == guest.Role {
			guests = append(guests, guest)
		}
	}

	if err = rows.Err(); err != nil {
		logs.LogError(err, "Get Guests Query Error")
		return guests, err
	}

	return guests, err
}

// RetrievePendingInvites opens a database connection and retrieves the list of guest users waiting for approval.
func RetrievePendingInvites(team string) ([]map[string]string, error) {
	var invites []map[string]string
	var err error
	var query string
	var rows *sql.Rows

	pool := data.ConnectToDB()
	defer pool.Close()

	if team == "" {
		query = `SELECT email, first_name, last_name, role, team, expiration, date_invited, proposer
			 FROM guests LEFT JOIN invites ON guests.email=invites.invitee
			 WHERE inviter IS NULL AND proposer IS NOT NULL AND pending=TRUE AND expiration >= NOW()
			 ORDER BY first_name;`
		rows, err = pool.Query(query)
	} else {
		query =
			`SELECT email, first_name, last_name, role, team, expiration, date_invited, proposer
			 FROM guests LEFT JOIN invites ON guests.email=invites.invitee
			 WHERE inviter IS NULL AND proposer IS NOT NULL AND pending=TRUE AND expiration >= NOW() AND team = $1
			 ORDER BY first_name;`
		rows, err = pool.Query(query, team)
	}

	if err != nil {
		logs.LogError(err, "Get Pending Invites Query Error")
		return invites, err
	}

	defer rows.Close()

	for rows.Next() {
		var guest data.GuestInvite
		if err := rows.Scan(&guest.Email, &guest.NameFirst, &guest.NameLast, &guest.Role, &guest.Team, &guest.Expires, &guest.DateInvited, &guest.Proposer); err != nil {
			logs.LogError(err, "Get Guests Query Error")
			return invites, err
		}

		guestData := map[string]string{
			"email":       guest.Email,
			"givenName":   guest.NameFirst,
			"familyName":  guest.NameLast,
			"role":        guest.Role,
			"team":        guest.Team,
			"expiration":  guest.Expires,
			"dateInvited": guest.DateInvited,
			"proposer":    guest.Proposer,
		}

		invites = append(invites, guestData)
	}

	if err = rows.Err(); err != nil {
		logs.LogError(err, "Get Pending Invites Query Error")
		return invites, err
	}

	return invites, err
}

// RetrieveUploaders returns a list of all guest users in a given team.
func RetrieveUploaders(team string) ([]map[string]any, error) {
	var uploaders []map[string]any

	pool := data.ConnectToDB()
	defer pool.Close()

	query :=
		`SELECT email, first_name, last_name, role, team, expiration, date_invited,
		 proposer, inviter, pending FROM guest_auth_data
		 WHERE team = $1 AND role = 'guest' ORDER BY first_name;`
	rows, err := pool.Query(query, team)

	if err != nil {
		logs.LogError(err, "Get Uploaders Query Error")
		return uploaders, err
	}

	defer rows.Close()

	for rows.Next() {
		var guest data.UploaderUser
		err := rows.Scan(
			&guest.Email,
			&guest.NameFirst,
			&guest.NameLast,
			&guest.Role,
			&guest.Team,
			&guest.Expires,
			&guest.DateInvited,
			&guest.Proposer,
			&guest.Inviter,
			&guest.Pending,
		)

		if err != nil {
			logs.LogError(err, "Get Uploaders Scan Error")
			return uploaders, err
		}

		guestData := map[string]any{
			"email":       guest.Email,
			"givenName":   guest.NameFirst,
			"familyName":  guest.NameLast,
			"role":        guest.Role,
			"team":        guest.Team,
			"expires":     guest.Expires,
			"dateInvited": guest.DateInvited,
			"proposer":    guest.Proposer.String,
			"inviter":     guest.Inviter.String,
			"pending":     guest.Pending,
		}

		uploaders = append(uploaders, guestData)
	}

	if err = rows.Err(); err != nil {
		logs.LogError(err, "Get Uploaders Row Error")
		return uploaders, err
	}

	return uploaders, err
}

// UpdateGuest opens a database connection and updates a given
// guest user with the provided information.
// TODO? - Allow for changes to user email? If so we may need
// to add an id field and set that as the primary key on a guest.
func UpdateGuest(guest data.GuestUser) error {
	pool := data.ConnectToDB()
	defer pool.Close()

	currentTime := time.Now()

	query :=
		`UPDATE guests SET first_name = $1, last_name = $2, role = $3,
		 team = $4, date_modified = $5 WHERE email = $6`
	_, err := pool.Exec(query, guest.NameFirst, guest.NameLast, guest.Role, guest.Team, currentTime, guest.Email)

	if err != nil {
		logs.LogError(err, "Update Guest Query Error")
	}

	return err
}

func shouldResetPassword(dateInvited string, nextExpiration time.Time, passwordWasReset bool) (bool, error) {
	var err error
	requireReset := true

	// Allow at most one reauthorization without a password reset
	if !passwordWasReset {
		return requireReset, err
	}

	parsedPrevInviteDate, err := time.Parse(time.RFC3339, dateInvited)

	if err != nil {
		return requireReset, err
	}

	// Reset after at most 60 days (use hours b/c no Days builtin)
	dateToReset := parsedPrevInviteDate.Add(time.Duration(60*24) * time.Hour)

	requireReset = dateToReset.Before(nextExpiration)

	return requireReset, err
}

func Reauthorize(guest data.GuestReauth, clientIsGuestAdmin bool) (string, int, error) {
	var pass string

	pool := data.ConnectToDB()
	defer pool.Close()

	var dateInvited string
	var pending bool
	var active bool
	var salt string
	var passHash string
	var passwordWasReset bool
	firstLogin := false

	query :=
		`SELECT date_invited, pending, expiration >= NOW() AS active, salt, pass_hash, password_reset
		 FROM invites WHERE invitee = $1 ORDER BY date_invited DESC LIMIT 1;`
	err := pool.QueryRow(query, guest.Email).Scan(&dateInvited, &pending, &active, &salt, &passHash, &passwordWasReset)

	if err != nil {
		return pass, 500, err
	} else if pending || active {
		return pass, 409, err
	}

	resetPassword, err := shouldResetPassword(dateInvited, guest.Expires, passwordWasReset)

	if err != nil {
		return pass, 500, err
	}

	if resetPassword {
		pass, salt = hashing.GenerateCredentials()
		passHash = hashing.GenerateHash(pass, salt)
		firstLogin = true
	}

	err = invites.SaveInvite(guest.Admin, guest.Email, guest.Expires, passHash, salt, clientIsGuestAdmin, resetPassword, firstLogin)

	return pass, 200, err
}

func AcceptGuest(guest data.AcceptInvite, hash string, salt string) error {
	pool := data.ConnectToDB()
	defer pool.Close()

	query := `UPDATE invites SET inviter = $1, pass_hash = $2, salt = $3, pending = FALSE WHERE invitee = $4`
	_, err := pool.Exec(query, guest.Inviter, hash, salt, guest.Invitee)

	if err != nil {
		logs.LogError(err, "Update Invite Query Error")
	}

	return err
}
