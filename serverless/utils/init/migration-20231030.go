package init

import (
	"database/sql"

	"github.com/IIP-Design/commons-gateway/utils/logs"
)

func createPasswordHistoryTable(pool *sql.DB) error {
	var err error

	_, err = pool.Exec(
		`CREATE TABLE IF NOT EXISTS password_history ( id VARCHAR(20) PRIMARY KEY, user_id VARCHAR(255) NOT NULL, creation_date TIMESTAMP NOT NULL, salt VARCHAR(10) NOT NULL, pass_hash VARCHAR(255) NOT NULL, FOREIGN KEY(user_id) REFERENCES guests(email) ON UPDATE CASCADE ON DELETE CASCADE );`,
	)

	if err != nil {
		logs.LogError(err, "Create Password History Table Query Error")
	}

	return err
}

// applyMigration20231030 prevents guests from reusing one of their last 24 passwords
func applyMigration20231030(title string) error {
	var err error

	pool := ConnectToDBInit()
	defer pool.Close()

	err = createPasswordHistoryTable(pool)

	if err != nil {
		return err
	}

	err = recordMigration(title)

	return err
}
