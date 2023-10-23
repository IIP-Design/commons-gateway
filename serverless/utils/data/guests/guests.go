package guests

import (
	"database/sql"
	"time"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
)

type InviteRecord struct {
	Pending     bool   `json:"pending"`
	DateInvited string `json:"dateInvited"`
	Expiration  string `json:"expiration"`
	Expired     bool   `json:"expired"`
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

	query = `SELECT pending, date_invited, expiration, expiration < NOW() AS expired FROM invites WHERE invitee = $1 ORDER BY expiration DESC`
	rows, err := pool.Query(query, email)

	if err != nil {
		logs.LogError(err, "Retrieve Invites Query Error")
		return guest, err
	}

	defer rows.Close()

	for rows.Next() {
		var pending bool
		var dateInvited time.Time
		var expiration time.Time
		var expired bool

		if err := rows.Scan(&pending, &dateInvited, &expiration, &expired); err != nil {
			logs.LogError(err, "Scan Guests Query Error")
			return guest, err
		}

		var invite = InviteRecord{
			Pending:     pending,
			DateInvited: dateInvited.Format(time.RFC3339),
			Expiration:  expiration.Format(time.RFC3339),
		}

		guest.Invites = append(guest.Invites, invite)
	}

	return guest, err
}

// RetrieveGuests opens a database connection and retrieves the full list of guest users.
func RetrieveGuests(team string, role string) ([]map[string]string, error) {
	var guests []map[string]string
	var err error
	var query string
	var rows *sql.Rows

	pool := data.ConnectToDB()
	defer pool.Close()

	if team == "" {
		query = `SELECT email, first_name, last_name, role, team, expiration FROM guest_auth_data ORDER BY first_name;`
		rows, err = pool.Query(query)
	} else {
		query =
			`SELECT email, first_name, last_name, role, team, expiration
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
		if err := rows.Scan(&guest.Email, &guest.NameFirst, &guest.NameLast, &guest.Role, &guest.Team, &guest.Expires); err != nil {
			logs.LogError(err, "Get Guests Query Error")
			return guests, err
		}

		guestData := map[string]string{
			"email":      guest.Email,
			"givenName":  guest.NameFirst,
			"familyName": guest.NameLast,
			"role":       guest.Role,
			"team":       guest.Team,
			"expiration": guest.Expires,
		}

		if role == "" || role == guestData["role"] {
			guests = append(guests, guestData)
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

func RetrieveUploaders(team string) ([]map[string]any, error) {
	var uploaders []map[string]any

	pool := data.ConnectToDB()
	defer pool.Close()

	query :=
		`SELECT email, first_name, last_name, role, team, expiration, date_invited, proposer, inviter, pending
			FROM guest_auth_data
			WHERE team = $1 ORDER BY first_name;`
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
			"expiration":  guest.Expires,
			"dateInvited": guest.DateInvited,
			"proposer":    guest.Proposer,
			"inviter":     guest.Inviter,
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
		`UPDATE guests SET first_name = $1, last_name = $2, team = $3,
		 date_modified = $4 WHERE email = $5`
	_, err := pool.Exec(query, guest.NameFirst, guest.NameLast, guest.Team, currentTime, guest.Email)

	if err != nil {
		logs.LogError(err, "Update Guest Query Error")
	}

	query =
		`UPDATE invites SET expiration = $1 WHERE invitee = $3`
	_, err = pool.Exec(query, guest.Expires, guest.Email)

	if err != nil {
		logs.LogError(err, "Update Invites Query Error")
	}

	return err
}

func AcceptGuest(guest data.AcceptInvite, hash string, salt string) error {
	pool := data.ConnectToDB()
	defer pool.Close()

	query :=
		`UPDATE invites SET inviter = $1, pass_hash = $2, salt = $3, pending = FALSE WHERE invitee = $4`
	_, err := pool.Exec(query, guest.Inviter, hash, salt, guest.Invitee)

	if err != nil {
		logs.LogError(err, "Update Invite Query Error")
	}

	return err
}
