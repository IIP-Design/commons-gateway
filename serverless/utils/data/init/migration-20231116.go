package init

import (
	"database/sql"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
)

// changeDescriptionType switches the type on upload descriptions from varchar(255) to text.
func changeDescriptionType(pool *sql.DB) error {
	var err error

	_, err = pool.Exec(
		`ALTER TABLE uploads ALTER COLUMN description TYPE TEXT;`,
	)

	if err != nil {
		logs.LogError(err, "Alter Description Column Type Query Error")
	}

	return err
}

// applyMigration20231116 removes the character limit on file descriptions.
func applyMigration20231116(title string) error {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	err = changeDescriptionType(pool)

	if err != nil {
		return err
	}

	err = recordMigration(title)

	return err
}
