package init

import (
	"database/sql"

	"github.com/IIP-Design/commons-gateway/utils/logs"
)

func movePassHashColumn(pool *sql.DB) error {
	var err error

	_, err = pool.Exec(
		`ALTER TABLE invites ADD COLUMN pass_hash VARCHAR(255) NOT NULL DEFAULT '';`,
	)

	if err != nil {
		logs.LogError(err, "Add Hash Column Query Error")
		return err
	}

	// Get all guest password hashes and add them to the invites table.
	guestRows, err := pool.Query(`SELECT email, pass_hash FROM guests;`)

	if err != nil {
		logs.LogError(err, "Select Guests Query Error")
		return err
	}

	defer guestRows.Close()

	for guestRows.Next() {
		var email string
		var passHash string

		if err := guestRows.Scan(&email, &passHash); err != nil {
			logs.LogError(err, "Scan Guests Query Error")
			return err
		}

		_, err = pool.Exec(`UPDATE invites SET pass_hash = $1 WHERE invitee = $2`, passHash, email)

		if err != nil {
			logs.LogError(err, "Invite pass_hash update error")

			return err
		}
	}

	_, err = pool.Exec(
		`ALTER TABLE guests DROP COLUMN pass_hash;`,
	)

	if err != nil {
		logs.LogError(err, "Drop Hash Column Query Error")
		return err
	}

	_, err = pool.Exec(
		`ALTER TABLE invites ALTER COLUMN pass_hash DROP DEFAULT;`,
	)

	if err != nil {
		logs.LogError(err, "Drop Hash Default Query Error")
		return err
	}

	return err
}

func moveSaltColumn(pool *sql.DB) error {
	var err error

	_, err = pool.Exec(
		`ALTER TABLE invites ADD COLUMN salt VARCHAR(10) NOT NULL DEFAULT '';`,
	)

	if err != nil {
		logs.LogError(err, "Add salt Column Query Error")
		return err
	}

	// Get all guest salts and add them to the invites table.
	guestRows, err := pool.Query(`SELECT email, salt FROM guests;`)

	if err != nil {
		logs.LogError(err, "Select Guests Query Error")
		return err
	}

	defer guestRows.Close()

	for guestRows.Next() {
		var email string
		var salt string

		if err := guestRows.Scan(&email, &salt); err != nil {
			logs.LogError(err, "Scan Guests Query Error")
			return err
		}

		_, err = pool.Exec(`UPDATE invites SET salt = $1 WHERE invitee = $2`, salt, email)

		if err != nil {
			logs.LogError(err, "Invite salt update error")

			return err
		}
	}

	_, err = pool.Exec(
		`ALTER TABLE guests DROP COLUMN salt;`,
	)

	if err != nil {
		logs.LogError(err, "Drop Salt Column Query Error")
		return err
	}

	_, err = pool.Exec(
		`ALTER TABLE invites ALTER COLUMN salt DROP DEFAULT;`,
	)

	if err != nil {
		logs.LogError(err, "Drop Salt Default Query Error")
		return err
	}

	return err
}

func moveExpirationColumn(pool *sql.DB) error {
	var err error

	_, err = pool.Exec(
		`ALTER TABLE invites ADD COLUMN expiration TIMESTAMP NOT NULL DEFAULT NOW();`,
	)

	if err != nil {
		logs.LogError(err, "Add Expiration Column Query Error")
		return err
	}

	// Get all guest salts and add them to the invites table.
	guestRows, err := pool.Query(`SELECT email, expiration FROM guests;`)

	if err != nil {
		logs.LogError(err, "Select Guests Query Error")
		return err
	}

	defer guestRows.Close()

	for guestRows.Next() {
		var email string
		var expiration string

		if err := guestRows.Scan(&email, &expiration); err != nil {
			logs.LogError(err, "Scan Guests Query Error")
			return err
		}

		_, err = pool.Exec(`UPDATE invites SET expiration = $1 WHERE invitee = $2`, expiration, email)

		if err != nil {
			logs.LogError(err, "Invite Expiration update error")

			return err
		}
	}

	_, err = pool.Exec(
		`ALTER TABLE guests DROP COLUMN expiration;`,
	)

	if err != nil {
		logs.LogError(err, "Drop Expiration Column Query Error")
		return err
	}

	_, err = pool.Exec(
		`ALTER TABLE invites ALTER COLUMN expiration DROP DEFAULT;`,
	)

	if err != nil {
		logs.LogError(err, "Drop Expiration Default Query Error")
		return err
	}

	return err
}

func moveFirstLoginColumn(pool *sql.DB) error {
	var err error

	_, err = pool.Exec(
		`ALTER TABLE invites ADD COLUMN first_login BOOLEAN NOT NULL DEFAULT true;`,
	)

	if err != nil {
		logs.LogError(err, "Add Login Column Query Error")
		return err
	}

	// Get all guest salts and add them to the invites table.
	guestRows, err := pool.Query(`SELECT email, first_login FROM guests;`)

	if err != nil {
		logs.LogError(err, "Select Guests Query Error")
		return err
	}

	defer guestRows.Close()

	for guestRows.Next() {
		var email string
		var firstLogin bool

		if err := guestRows.Scan(&email, &firstLogin); err != nil {
			logs.LogError(err, "Scan Guests Query Error")
			return err
		}

		_, err = pool.Exec(`UPDATE invites SET first_login = $1 WHERE invitee = $2`, firstLogin, email)

		if err != nil {
			logs.LogError(err, "Invite Login update error")

			return err
		}
	}

	_, err = pool.Exec(
		`ALTER TABLE guests DROP COLUMN first_login;`,
	)

	if err != nil {
		logs.LogError(err, "Drop Login Column Query Error")
		return err
	}

	return err
}

func createInviteView(pool *sql.DB) error {
	var err error

	_, err = pool.Exec(
		`CREATE OR REPLACE VIEW recent_invites AS ( SELECT DISTINCT ON(invitee) * FROM invites ORDER BY invitee, date_invited DESC );`,
	)

	if err != nil {
		logs.LogError(err, "Create View Query Error")
	}

	return err
}

func createAuthView(pool *sql.DB) error {
	var err error

	_, err = pool.Exec(
		`CREATE OR REPLACE VIEW guest_auth_data AS ( SELECT * FROM guests LEFT JOIN recent_invites ON guests.email = recent_invites.invitee );`,
	)

	if err != nil {
		logs.LogError(err, "Create View Query Error")
	}

	return err
}

// applyMigration20231023 moves access data from guests table to invites table
func applyMigration20231023(title string) error {
	var err error

	pool := connectToDB()
	defer pool.Close()

	err = movePassHashColumn(pool)

	if err != nil {
		return err
	}

	err = moveSaltColumn(pool)

	if err != nil {
		return err
	}

	err = moveExpirationColumn(pool)

	if err != nil {
		return err
	}

	err = moveFirstLoginColumn(pool)

	if err != nil {
		return err
	}

	err = createInviteView(pool)

	if err != nil {
		return err
	}

	err = createAuthView(pool)

	if err != nil {
		return err
	}

	err = recordMigration(title)

	return err
}
