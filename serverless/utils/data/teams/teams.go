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
func CreateTeam(teamName string, aprimoName string) error {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	guid := xid.New()
	currentTime := time.Now()

	insertTeam :=
		`INSERT INTO teams( id, team_name, aprimo_name, active, date_created, date_modified )
		 VALUES ($1, $2, $3, $4, $5, $6);`
	_, err = pool.Exec(insertTeam, guid, teamName, aprimoName, true, currentTime, currentTime)

	if err != nil {
		logs.LogError(err, "Create Team Query Error")
	}

	return err
}

// GetTeamIdByName uses a team's name to retrieve it's unique identifier.
func GetTeamIdByName(teamName string) (string, error) {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	var id string

	query := `SELECT id FROM teams WHERE team_name = $1;`
	err = pool.QueryRow(query, teamName).Scan(&id)

	if err != nil {
		logs.LogError(err, "Existing Team by ID Query Error")
		return "", err
	}

	return id, err
}

// UpdateTeam opens a database connection and updates and existing team record.
func UpdateTeam(teamId string, teamName string, aprimoName string, active bool) error {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	currentTime := time.Now()

	query := `UPDATE teams SET team_name = $1, aprimo_name = $2, active = $3, date_modified = $4 WHERE id = $5;`
	_, err = pool.Exec(query, teamName, aprimoName, active, currentTime, teamId)

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

	currentTime := time.Now()

	query := `UPDATE teams SET active = $1, date_modified = $2 WHERE id = $3;`
	_, err = pool.Exec(query, active, currentTime, teamId)

	if err != nil {
		logs.LogError(err, "Update Team Status Query Error")
	}

	return err
}

// RetrieveTeams opens a database connection and retrieves the full list of teams.
func RetrieveTeams() ([]data.Team, error) {
	var teams []data.Team
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	rows, err := pool.Query(`SELECT id, team_name, active, aprimo_name FROM teams ORDER BY team_name`)

	if err != nil {
		logs.LogError(err, "Get Teams Query Error")
		return teams, err
	}

	defer rows.Close()

	for rows.Next() {
		var team data.Team
		if err := rows.Scan(&team.Id, &team.Name, &team.Active, &team.AprimoName); err != nil {
			logs.LogError(err, "Get Teams Query Error")
			return teams, err
		}

		teams = append(teams, team)
	}

	if err = rows.Err(); err != nil {
		logs.LogError(err, "Get Teams Query Error")
		return teams, err
	}

	return teams, err
}
