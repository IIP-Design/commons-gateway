package teams

import (
	"database/sql"
	"time"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"

	"github.com/rs/xid"
)

// CheckForExistingTeam opens a database connection and checks whether the provided team name
// is present. An affirmative check indicates that the given team has already been added.
func CheckForExistingTeam(teamName string) (bool, error) {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	var team string

	query := `SELECT team_name FROM teams WHERE team_name = $1;`
	err = pool.QueryRow(query, teamName).Scan(&team)

	if err != nil {
		// Do not return an error if no results are found.
		if err == sql.ErrNoRows {
			return false, nil
		}

		logs.LogError(err, "Existing Team Query Error")
	}

	return team == teamName, err
}

// CheckForExistingTeamById opens a database connection and checks whether the provided
// team id (which is a unique value constraint in the teams tables) is present.
// An affirmative check indicates that the given team has already been added.
func CheckForExistingTeamById(teamId string) (bool, error) {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	var team string

	query := `SELECT id FROM teams WHERE id = $1;`
	err = pool.QueryRow(query, teamId).Scan(&team)

	if err != nil {
		// Do not return an error if no results are found.
		if err == sql.ErrNoRows {
			return false, nil
		}

		logs.LogError(err, "Existing Team by ID Query Error")
	}

	return team == teamId, err
}

// CreateTeam opens a database connection and saves a new team record.
func CreateTeam(teamName string) error {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	guid := xid.New()
	currentTime := time.Now()

	insertTeam :=
		`INSERT INTO "teams"("id", "team_name", "active", "date_created")
		 VALUES ($1, $2, $3, $4);`
	_, err = pool.Exec(insertTeam, guid, teamName, true, currentTime)

	if err != nil {
		logs.LogError(err, "Create Team Query Error")
	}

	return err
}

// UpdateTeam opens a database connection and updates and existing team record.
func UpdateTeam(teamId string, teamName string, active bool) error {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	query := `UPDATE teams SET team_name = $1, active = $2 WHERE id = $3;`
	_, err = pool.Exec(query, teamName, active, teamId)

	if err != nil {
		logs.LogError(err, "Update Team Query Error")
	}

	return err
}

// UpdateTeamStatus opens a database connection and updates and existing team's status.
func UpdateTeamStatus(teamId string, active bool) error {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	query := `UPDATE teams SET active = $1 WHERE id = $2;`
	_, err = pool.Exec(query, active, teamId)

	if err != nil {
		logs.LogError(err, "Update Team Status Query Error")
	}

	return err
}

// RetrieveTeams opens a database connection and retrieves the full list of teams.
func RetrieveTeams() ([]map[string]any, error) {
	var teams []map[string]any
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	rows, err := pool.Query(`SELECT id, team_name, active FROM teams ORDER BY team_name`)

	if err != nil {
		logs.LogError(err, "Get Teams Query Error")
		return teams, err
	}

	defer rows.Close()

	for rows.Next() {
		var team data.Team
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
