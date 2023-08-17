package data

import (
	"encoding/json"

	"github.com/IIP-Design/commons-gateway/utils/logs"
)

// RetrieveTeams opens a database connection and retrieves the full list of teams.
func RetrieveTeams() ([]string, error) {
	var teams []string
	var err error

	pool := connectToDB()
	defer pool.Close()

	rows, err := pool.Query(`SELECT id, team_name FROM teams`)

	if err != nil {
		logs.LogError(err, "Get Teams Query Error")
		return teams, err
	}

	defer rows.Close()

	for rows.Next() {
		var team Team
		if err := rows.Scan(&team.Id, &team.Name); err != nil {
			logs.LogError(err, "Get Teams Query Error")
			return teams, err
		}

		bytes, err := json.Marshal(team)

		if err != nil {
			logs.LogError(err, "Failed to Marshal Team User Data.")
			return teams, err
		}

		teams = append(teams, string(bytes[:]))
	}

	if err = rows.Err(); err != nil {
		logs.LogError(err, "Get Teams Query Error")
		return teams, err
	}

	return teams, err
}
