package data

import (
	"encoding/json"
	"fmt"

	"github.com/IIP-Design/commons-gateway/utils/logs"
)

// RetrieveGuests opens a database connection and retrieves the full list of admin users.
func RetrieveGuests(team string) ([]string, error) {
	var guests []string
	var err error
	var query string

	pool := connectToDB()
	defer pool.Close()

	if team == "" {
		query = `SELECT email, first_name, last_name, team, expiration FROM guests`
	} else {
		query = fmt.Sprintf(
			`SELECT email, first_name, last_name, team, expiration FROM guests WHERE team = '%s';`,
			team,
		)
	}

	rows, err := pool.Query(query)

	if err != nil {
		logs.LogError(err, "Get Guests Query Error")
		return guests, err
	}

	defer rows.Close()

	for rows.Next() {
		var guest GuestUser
		if err := rows.Scan(&guest.Email, &guest.NameFirst, &guest.NameLast, &guest.Team, &guest.Expires); err != nil {
			logs.LogError(err, "Get Guests Query Error")
			return guests, err
		}

		bytes, err := json.Marshal(guest)

		if err != nil {
			logs.LogError(err, "Failed to Marshal Guest User Data.")
			return guests, err
		}

		guests = append(guests, string(bytes[:]))
	}

	if err = rows.Err(); err != nil {
		logs.LogError(err, "Get Guests Query Error")
		return guests, err
	}

	return guests, err
}
