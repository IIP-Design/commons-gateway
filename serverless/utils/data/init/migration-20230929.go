package init

import (
	"database/sql"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
)

func addAprimoRecordId(pool *sql.DB) error {
	var err error

	_, err = pool.Exec(
		`ALTER TABLE uploads ADD COLUMN aprimo_record_id VARCHAR(255) DEFAULT NULL;`,
	)

	if err != nil {
		logs.LogError(err, "Alter Table Query Error - aprimo_record_id")

		return err
	}

	return err
}

// applyMigration20230929 adds support for Aprimo IDs to be stored with uploads.
func applyMigration20230929(title string) error {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	err = addAprimoRecordId(pool)

	if err != nil {
		return err
	}

	err = recordMigration(title)

	return err
}
