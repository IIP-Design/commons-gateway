package init

import (
	"database/sql"

	"github.com/IIP-Design/commons-gateway/utils/logs"
)

func addInvitePasswordResetColumn(pool *sql.DB) error {
	var err error

	_, err = pool.Exec(
		`ALTER TABLE invites ADD COLUMN password_reset BOOLEAN NOT NULL DEFAULT TRUE;`,
	)

	if err != nil {
		logs.LogError(err, "Add Hash Column Query Error")
	}

	return err
}

// applyMigration20231024 allows guests to skip resetting their password on certain reauthorizations
func applyMigration20231024(title string) error {
	var err error

	pool := connectToDB()
	defer pool.Close()

	err = addInvitePasswordResetColumn(pool)

	if err != nil {
		return err
	}

	err = recordMigration(title)

	return err
}
