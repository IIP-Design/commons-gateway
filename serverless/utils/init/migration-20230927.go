package init

import (
	"database/sql"

	"github.com/IIP-Design/commons-gateway/utils/logs"
)

// create2FATable adds a table to the database which stores an
// ephemeral list of second factor authentication requests.
func create2FATable(pool *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS mfa (
		request_id VARCHAR(20) PRIMARY KEY,
		code VARCHAR(20) NOT NULL,
		date_created TIMESTAMP NOT NULL
	);`

	_, err := pool.Exec(query)

	if err != nil {
		logs.LogError(err, "Table Creation Query Error - MFA")
	}

	return err
}

// applyMigration20230927 adds support for second factor
// authentication for partner user logins.
func applyMigration20230927(title string) error {
	var err error

	pool := connectToDB()
	defer pool.Close()

	err = create2FATable(pool)

	if err != nil {
		return err
	}

	err = recordMigration(title)

	return err
}
