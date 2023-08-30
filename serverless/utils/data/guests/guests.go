package guests

import (
	"database/sql"
	"time"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
)

// RetrieveGuest opens a database connection and retrieves the information for a single user.
func RetrieveGuest(email string) (map[string]string, error) {
	var guest map[string]string

	pool := data.ConnectToDB()
	defer pool.Close()

	var firstName string
	var lastName string
	var team string
	var expiration time.Time

	query := `SELECT first_name, last_name, team, expiration FROM guests WHERE email = $1`
	err := pool.QueryRow(query, email).Scan(&firstName, &lastName, &team, &expiration)

	if err != nil {
		logs.LogError(err, "Retrieve Guest Query Error")
	}

	guest = map[string]string{
		"email":      email,
		"givenName":  firstName,
		"familyName": lastName,
		"team":       team,
		"expiration": expiration.String(),
	}

	return guest, err
}

// RetrieveGuests opens a database connection and retrieves the full list of admin users.
func RetrieveGuests(team string) ([]map[string]string, error) {
	var guests []map[string]string
	var err error
	var query string
	var rows *sql.Rows

	pool := data.ConnectToDB()
	defer pool.Close()

	if team == "" {
		query = `SELECT email, first_name, last_name, team, expiration FROM guests ORDER BY first_name;`
		rows, err = pool.Query(query)
	} else {
		query =
			`SELECT email, first_name, last_name, team, expiration
			 FROM guests WHERE team = $1 ORDER BY first_name;`
		rows, err = pool.Query(query, team)
	}

	if err != nil {
		logs.LogError(err, "Get Guests Query Error")
		return guests, err
	}

	defer rows.Close()

	for rows.Next() {
		var guest data.GuestUser
		if err := rows.Scan(&guest.Email, &guest.NameFirst, &guest.NameLast, &guest.Team, &guest.Expires); err != nil {
			logs.LogError(err, "Get Guests Query Error")
			return guests, err
		}

		guestData := map[string]string{
			"email":      guest.Email,
			"givenName":  guest.NameFirst,
			"familyName": guest.NameLast,
			"team":       guest.Team,
			"expiration": guest.Expires,
		}

		guests = append(guests, guestData)
	}

	if err = rows.Err(); err != nil {
		logs.LogError(err, "Get Guests Query Error")
		return guests, err
	}

	return guests, err
}

// UpdateGuest opens a database connection and updates a given
// guest user with the provided information.
// TODO? - Allow for changes to user email? If so we may need
// to add an id field and set that as the primary key on a guest.
func UpdateGuest(guest data.GuestUser) error {
	pool := data.ConnectToDB()
	defer pool.Close()

	query :=
		`UPDATE guests SET first_name = $1, last_name = $2, team = $3,
		 expiration = $4 WHERE email = $5`
	_, err := pool.Exec(query, guest.NameFirst, guest.NameLast, guest.Team, guest.Expires, guest.Email)

	if err != nil {
		logs.LogError(err, "Update Guest Query Error")
	}

	return err
}
