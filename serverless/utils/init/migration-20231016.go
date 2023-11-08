package init

import (
	"database/sql"

	"github.com/IIP-Design/commons-gateway/utils/logs"
)

func updateS3KeyColumn(pool *sql.DB) error {
	var err error

	_, err = pool.Exec(
		`ALTER TABLE uploads ALTER COLUMN s3_id TYPE VARCHAR(255);`,
	)

	if err != nil {
		logs.LogError(err, "Update S3 Key Column Query Error")
	}

	return err
}

// applyMigration20231016 increases the size of the s3_id property since we no
// longer use a defined length file name.
func applyMigration20231016(title string) error {
	var err error

	pool := ConnectToDBInit()
	defer pool.Close()

	err = updateS3KeyColumn(pool)

	if err != nil {
		return err
	}

	err = recordMigration(title)

	return err
}
