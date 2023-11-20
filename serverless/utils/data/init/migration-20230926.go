package init

import (
	"database/sql"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
)

// addAprimoNameColumn adds a column to the teams table to keep track of team
// name equivalents used in the Aprimo system. It sets the initial value one
// existing rows to that team's name - which may not be the actual Aprimo name.
func addAprimoNameColumn(pool *sql.DB) error {
	var err error

	// Add new column to teams table.
	_, err = pool.Exec(
		`ALTER TABLE teams ADD COLUMN aprimo_name VARCHAR(255);`,
	)

	if err != nil {
		logs.LogError(err, "Add Aprimo Team Name Column Query Error")

		return err
	}

	// Copy over the team name as the Aprimo name. Note that this will
	// not be the actual Aprimo name, which will need to be manually set.
	_, err = pool.Exec(
		`UPDATE teams SET aprimo_name = team_name;`,
	)

	if err != nil {
		logs.LogError(err, "Populate Aprimo Team Name Column Query Error")

		return err
	}

	// Now that it is populated, add a not null constraint to the column.
	_, err = pool.Exec(
		`ALTER TABLE teams ALTER COLUMN aprimo_name SET NOT NULL;`,
	)

	if err != nil {
		logs.LogError(err, "Set Aprimo Team Name Column Null Constraint Query Error")

		return err
	}

	return err
}

// applyMigration20230926 enables better integration with Aprimo
// by supplying team names that match those in Aprimo.
func applyMigration20230926(title string) error {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	err = addAprimoNameColumn(pool)

	if err != nil {
		return err
	}

	err = recordMigration(title)

	return err
}
