package data

import (
	"github.com/IIP-Design/commons-gateway/utils/logs"
)

// RetrieveTeams opens a database connection and retrieves the full list of teams.
func RetrieveTeams() ([]map[string]any, error) {
	var teams []map[string]any
	var err error

	pool := connectToDB()
	defer pool.Close()

	rows, err := pool.Query(`SELECT id, team_name, active FROM teams`)

	if err != nil {
		logs.LogError(err, "Get Teams Query Error")
		return teams, err
	}

	defer rows.Close()

	for rows.Next() {
		var team Team
		if err := rows.Scan(&team.Id, &team.Name, &team.Active); err != nil {
			logs.LogError(err, "Get Teams Query Error")
			return teams, err
		}

		teamData := map[string]any{
			"id":     team.Id,
			"name":   team.Name,
			"active": team.Active,
		}

		teams = append(teams, teamData)
	}

	if err = rows.Err(); err != nil {
		logs.LogError(err, "Get Teams Query Error")
		return teams, err
	}

	return teams, err
}
