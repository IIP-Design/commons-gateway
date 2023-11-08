package init

import (
	"database/sql"

	"github.com/IIP-Design/commons-gateway/utils/logs"
)

// addLoginColumns adds columns to the guests table to track a given
// user's login attempts and (if relevant their lock out status).
func addLoginColumns(pool *sql.DB) error {
	var err error

	_, err = pool.Exec(
		`ALTER TABLE guests
		 ADD COLUMN first_login BOOLEAN NOT NULL DEFAULT true,
		 ADD COLUMN locked BOOLEAN NOT NULL DEFAULT false,
		 ADD COLUMN login_attempt SMALLINT DEFAULT 0,
		 ADD COLUMN login_date TIMESTAMP;`,
	)

	if err != nil {
		logs.LogError(err, "Add Guest Login Fields Query Error")
	}

	return err
}

// applyMigration20231010 tracks partner user login attempts and lockout status.
func applyMigration20231010(title string) error {
	var err error

	pool := ConnectToDBInit()
	defer pool.Close()

	err = addLoginColumns(pool)

	if err != nil {
		return err
	}

	err = recordMigration(title)

	return err
}
