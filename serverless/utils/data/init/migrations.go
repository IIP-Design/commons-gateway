package init

import (
	"fmt"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
)

// A list of schema migrations. Please add newest to the bottom.
// They should also follow the naming conventions `migYYYYMMDD` for
// the const name and `YYYYMMDD_short_description` for the value.
const mig20230831 = "20230831_user_roles"
const mig20230926 = "20230926_aprimo_teams"
const mig20230927 = "20230927_mfa_table"
const mig20230929 = "20230929_aprimo_ids"
const mig20231002 = "20231002_aprimo_tokens"

// getAppliedMigrations queries the `migrations` table in that database
// for a list of schema updates that have already been executed.
func getAppliedMigrations() ([]string, error) {
	var applied []string
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	rows, err := pool.Query(`SELECT title FROM migrations`)

	if err != nil {
		logs.LogError(err, "Fetch Migrations Error")
		return applied, err
	}

	defer rows.Close()

	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			logs.LogError(err, "Get Migrations Title Error")
			return applied, err
		}

		applied = append(applied, title)
	}

	return applied, err
}

// ApplyMigrations loops through the list of schema updates already applied to the
// database. If an update has not already been applied, it is executed. Migrations
// should be listed in chronological order, with the most recent appearing last.
func ApplyMigrations() error {
	var err error

	applied, err := getAppliedMigrations()

	if err != nil {
		return err
	}

	// Apply the migration from August 31, 2023
	if !stringArrayContains(applied, mig20230831) {
		fmt.Printf("Applying migration - %s", mig20230831)

		err = applyMigration20230831(mig20230831)

		if err != nil {
			return err
		}
	}

	// Apply the migration from September 26, 2023
	if !stringArrayContains(applied, mig20230926) {
		fmt.Printf("Applying migration - %s", mig20230926)

		err = applyMigration20230926(mig20230926)

		if err != nil {
			return err
		}
	}

	// Apply the migration from September 27, 2023
	if !stringArrayContains(applied, mig20230927) {
		fmt.Printf("Applying migration - %s", mig20230927)

		err = applyMigration20230927(mig20230927)

		if err != nil {
			return err
		}
	}

	// Apply the migration from September 29, 2023
	if !stringArrayContains(applied, mig20230929) {
		fmt.Printf("Applying migration - %s", mig20230929)

		err = applyMigration20230929(mig20230929)

		if err != nil {
			return err
		}
	}

	// Apply the migration from October 02, 2023
	if !stringArrayContains(applied, mig20231002) {
		fmt.Printf("Applying migration - %s", mig20231002)

		err = applyMigration20231002(mig20231002)

		if err != nil {
			return err
		}
	}

	return err
}
