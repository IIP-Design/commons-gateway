package init

import (
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
)

// InitializeDatabase opens a connection to the database and runs
// the queries needed to configure the database with the proper tables.
func InitializeDatabase() error {
	var err error
	pool := data.ConnectToDB()
	defer pool.Close()

	_, err = pool.Exec(teamsQuery)

	if err != nil {
		logs.LogError(err, "Table Creation Error - Teams")
	}

	_, err = pool.Exec(adminsQuery)

	if err != nil {
		logs.LogError(err, "Table Creation Error - Admins")
	}

	_, err = pool.Exec(guestsQuery)

	if err != nil {
		logs.LogError(err, "Table Creation Error - Guests")
	}

	_, err = pool.Exec(invitesQuery)

	if err != nil {
		logs.LogError(err, "Table Creation Error - Invites")
	}

	_, err = pool.Exec(uploadsQuery)

	if err != nil {
		logs.LogError(err, "Table Creation Error - Uploads")
	}

	_, err = pool.Exec(migrationsQuery)

	if err != nil {
		logs.LogError(err, "Table Creation Error - Migrations")
	}

	return err
}
