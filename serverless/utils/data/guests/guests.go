package guests

import (
	"fmt"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
)

// RetrieveGuests opens a database connection and retrieves the full list of admin users.
func RetrieveGuests(team string) ([]map[string]string, error) {
	var guests []map[string]string
	var err error
	var query string

	pool := data.ConnectToDB()
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
