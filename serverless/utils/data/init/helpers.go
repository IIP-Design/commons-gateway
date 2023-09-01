package init

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/rs/xid"
)

// CheckForTable identifies whether or not a table with a given name
// is present in the database.
func CheckForTable(tablename string) bool {
	var exists bool

	pool := data.ConnectToDB()
	defer pool.Close()

	query :=
		`SELECT EXISTS ( SELECT FROM pg_tables WHERE schemaname = 'public'
		 AND tablename  = $1 );`

	err := pool.QueryRow(query, tablename).Scan(&exists)

	if err == sql.ErrNoRows {
		msg := fmt.Sprintf("Table Not Found - %s", tablename)

		logs.LogError(err, msg)
		return false
	} else if err != nil {
		msg := fmt.Sprintf("Check Table Query Error - %s", tablename)

		logs.LogError(err, msg)
		return false
	}

	return exists
}

// stringArrayContains iterates through an array of string to check whether
// a given string value is present as an item in the array.
func stringArrayContains(arr []string, value string) bool {
	for _, item := range arr {
		if item == value {
			return true
		}
	}

	return false
}

// recordMigration saves the title of an schema migration and the date on which it was applied.
func recordMigration(title string) error {
	pool := data.ConnectToDB()
	defer pool.Close()

	guid := xid.New()
	currentTime := time.Now()

	query := `INSERT INTO migrations( id, title, date_applied ) VALUES ( $1, $2, $3 );`

	_, err := pool.Exec(query, guid, title, currentTime)

	if err != nil {
		logs.LogError(err, "Migration Registry Error")
	}

	return err
}
