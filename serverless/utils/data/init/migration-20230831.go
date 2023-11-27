package init

import (
	"database/sql"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/data/users"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/rs/xid"
)

// createAllUsersTable adds a table to the database which stores a list
// of all admin and guest users.
func createAllUsersTable(pool *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS all_users (
		user_id VARCHAR(20) PRIMARY KEY,
		admin_id VARCHAR(255),
		guest_id VARCHAR(255),
		FOREIGN KEY(admin_id) REFERENCES admins(email) ON UPDATE CASCADE ON DELETE CASCADE,
		FOREIGN KEY(guest_id) REFERENCES guests(email) ON UPDATE CASCADE ON DELETE CASCADE
	);`

	_, err := pool.Exec(query)

	if err != nil {
		logs.LogError(err, "Table Creation Query Error - All Users")
	}

	return err
}

// populateAllUsersTable iterates through all existing admin and guest users
// adding them to the `all_users` table.
func populateAllUsersTable(pool *sql.DB) error {
	// Get all guests and add them to the users table.
	guestRows, err := pool.Query(`SELECT email FROM guests;`)

	if err != nil {
		logs.LogError(err, "Select Guests Query Error")
		return err
	}

	defer guestRows.Close()

	for guestRows.Next() {
		var email string
		if err := guestRows.Scan(&email); err != nil {
			logs.LogError(err, "Select Guests Query Error")
			return err
		}

		guid := xid.New()

		_, err = pool.Exec(`INSERT INTO all_users( user_id, guest_id ) VALUES ( $1, $2 )`, guid, email)

		if err != nil {
			logs.LogError(err, "All User Insert Query Error")

			return err
		}
	}

	// Get all admins and add them to the users table.
	adminRows, err := pool.Query(`SELECT email FROM admins;`)

	if err != nil {
		logs.LogError(err, "Select Admins Query Error")
		return err
	}

	defer adminRows.Close()

	for adminRows.Next() {
		var email string
		if err := adminRows.Scan(&email); err != nil {
			logs.LogError(err, "Select Admin Query Error")
			return err
		}

		guid := xid.New()

		_, err = pool.Exec(`INSERT INTO all_users( user_id, admin_id ) VALUES ( $1, $2 )`, guid, email)

		if err != nil {
			logs.LogError(err, "All User Insert Query Error")

			return err
		}
	}

	return err
}

// switchUploadsUserId changes the `user_id` value for any uploads from the user's email
// from either the `admins` or `guests` table to their `user_id` from the `all_users` table
func switchUploadsUserId(pool *sql.DB) error {
	// Get all the user_ids from the uploads table
	uploadsRows, err := pool.Query(`SELECT s3_id, user_id FROM uploads;`)

	if err != nil {
		logs.LogError(err, "Select Uploads Query Error")
		return err
	}

	defer uploadsRows.Close()

	// Loop through the user_ids updating the uploads
	for uploadsRows.Next() {
		var uploadId string
		var oldId string
		if err := uploadsRows.Scan(&uploadId, &oldId); err != nil {
			logs.LogError(err, "Select Upload User Id Query Error")
			return err
		}

		// Check if user exists in the admins table
		_, isAdmin, err := users.CheckForExistingAdminUser(oldId)

		if err == nil && isAdmin {
			var userId string

			// Retrieve the new user id from the all_users table
			err := pool.QueryRow(`SELECT user_id FROM all_users WHERE admin_id = $1`, oldId).Scan(&userId)

			if err != nil {
				logs.LogError(err, "Select Admin's User Id Query Error")

				return err
			}

			_, err = pool.Exec(
				`UPDATE uploads SET user_id = $1 WHERE s3_id = $2;`,
				userId,
				uploadId,
			)

			if err != nil {
				logs.LogError(err, "Update Upload User Id Query Error")

				return err
			}
		} else if err == nil && !isAdmin {
			var userId string

			// If not found in admin table, look in the guests table.
			_, isGuest, err := users.CheckForExistingGuestUser(oldId)

			if err != nil {
				logs.LogError(err, "Check for Guest Query Error")

				return err
			} else if isGuest {
				// Retrieve the new user id from the all_users table
				err := pool.QueryRow(`SELECT user_id FROM all_users WHERE guest_id = $1`, oldId).Scan(&userId)

				if err != nil {
					logs.LogError(err, "Select Guest's User Id Query Error")

					return err
				}

				_, err = pool.Exec(
					`UPDATE uploads SET user_id = $1 WHERE s3_id = $2;`,
					userId,
					uploadId,
				)

				if err != nil {
					logs.LogError(err, "Update Upload User Id Query Error")

					return err
				}
			}
		} else if err != nil {
			logs.LogError(err, "Check for Admin User Query Error")

			return err
		}
	}

	// Set foreign key constraint on `user_id` property.
	_, err = pool.Exec(
		`ALTER TABLE uploads ADD CONSTRAINT uploads_user_id_fkey FOREIGN KEY(user_id)
		 REFERENCES all_users(user_id) ON UPDATE CASCADE ON DELETE RESTRICT;`,
	)

	if err != nil {
		logs.LogError(err, "Update Uploads User Id FK Query Error")
	}

	return err
}

// enableUserRoles enumerates two possible admin roles and two possible
// guest roles. It then adds a column to the `admins` and `guests` tables
// so that each user can be assigned one of those roles.
func enableUserRoles(pool *sql.DB) error {
	var err error

	// Create and set admin roles.
	_, err = pool.Exec(`CREATE TYPE ADMIN_ROLE AS ENUM ( 'super admin', 'admin' );`)

	if err != nil {
		logs.LogError(err, "Add Admin Role Enum Query Error")

		return err
	}

	_, err = pool.Exec(`ALTER TABLE admins ADD COLUMN role ADMIN_ROLE NOT NULL DEFAULT 'admin';`)

	if err != nil {
		logs.LogError(err, "Add Role to Admins Table Query Error")

		return err
	}

	// Create and set guest roles.
	_, err = pool.Exec(`CREATE TYPE GUEST_ROLE AS ENUM ( 'guest admin', 'guest' );`)

	if err != nil {
		logs.LogError(err, "Add Guest Role Enum Query Error")

		return err
	}

	_, err = pool.Exec(`ALTER TABLE guests ADD COLUMN role GUEST_ROLE NOT NULL DEFAULT 'guest';`)

	if err != nil {
		logs.LogError(err, "Add Role to Guests Query Error")

		return err
	}

	return err
}

// addInviteProposerSupport adds optional fields to the invite table that are required
// if/when guest admins are allowed to propose invitations.
func addInviteProposerSupport(pool *sql.DB) error {
	var err error

	_, err = pool.Exec(
		`ALTER TABLE invites
		 ADD COLUMN pending BOOLEAN NOT NULL DEFAULT false,
		 ADD COLUMN proposer VARCHAR(255) REFERENCES guests(email)
		 ON UPDATE CASCADE ON DELETE RESTRICT;`,
	)

	if err != nil {
		logs.LogError(err, "Add Proposed Invite Fields Query Error")

		return err
	}

	// Make inviter nullable - for when pending admin approval.
	_, err = pool.Exec(`ALTER TABLE invites ALTER COLUMN inviter DROP NOT NULL;`)

	if err != nil {
		logs.LogError(err, "Drop Inviter Non-Null Query Error")
	}

	return err
}

// addModifiedDates creates a new column for admins, guests, and users to indicate
// when they were last modified. We set the initial value to the date of the migration
// only so that we can add a non-null constraint without errors
func addModifiedDates(pool *sql.DB) error {
	var err error

	// Uploads
	_, err = pool.Exec(
		`ALTER TABLE uploads ADD COLUMN date_uploaded TIMESTAMP NOT NULL DEFAULT current_timestamp;`,
	)

	if err != nil {
		logs.LogError(err, "Add Uploaded Date to Uploads Query Error")

		return err
	}

	_, err = pool.Exec(
		`ALTER TABLE uploads ALTER COLUMN date_uploaded DROP DEFAULT;`,
	)

	if err != nil {
		logs.LogError(err, "Drop Uploads Uploaded Default Query Error")

		return err
	}

	// Teams
	_, err = pool.Exec(
		`ALTER TABLE teams ADD COLUMN date_modified TIMESTAMP NOT NULL DEFAULT current_timestamp;`,
	)

	if err != nil {
		logs.LogError(err, "Add Modified Date to Teams Query Error")

		return err
	}

	_, err = pool.Exec(
		`ALTER TABLE teams ALTER COLUMN date_modified DROP DEFAULT;`,
	)

	if err != nil {
		logs.LogError(err, "Drop Teams Modified Default Query Error")

		return err
	}

	// Guests
	_, err = pool.Exec(
		`ALTER TABLE guests ADD COLUMN date_modified TIMESTAMP NOT NULL DEFAULT current_timestamp;`,
	)

	if err != nil {
		logs.LogError(err, "Add Modified Date to Guests Query Error")

		return err
	}

	_, err = pool.Exec(
		`ALTER TABLE guests ALTER COLUMN date_modified DROP DEFAULT;`,
	)

	if err != nil {
		logs.LogError(err, "Drop Guests Modified Default Query Error")

		return err
	}

	// Admins
	_, err = pool.Exec(
		`ALTER TABLE admins ADD COLUMN date_modified TIMESTAMP NOT NULL DEFAULT current_timestamp;`,
	)

	if err != nil {
		logs.LogError(err, "Add Modified Date to Admins Query Error")

		return err
	}

	_, err = pool.Exec(
		`ALTER TABLE admins ALTER COLUMN date_modified DROP DEFAULT;`,
	)

	if err != nil {
		logs.LogError(err, "Drop Admin Modified Default Query Error")
	}

	return err
}

// updateConstraints switches some cascade on deletes to restrict on deletes.
// This is done out of an abundance of caution to ensure data integrity.
func updateConstraints(pool *sql.DB) error {
	var err error

	_, err = pool.Exec(
		`ALTER TABLE admins DROP CONSTRAINT admins_team_fkey,
		 ADD CONSTRAINT admins_team_fkey FOREIGN KEY(team)
		 REFERENCES teams(id) ON UPDATE CASCADE ON DELETE RESTRICT;`,
	)

	if err != nil {
		logs.LogError(err, "Update Admins Team FK Query Error")

		return err
	}

	_, err = pool.Exec(
		`ALTER TABLE guests DROP CONSTRAINT guests_team_fkey,
		 ADD CONSTRAINT guests_team_fkey FOREIGN KEY(team)
		 REFERENCES teams(id) ON UPDATE CASCADE ON DELETE RESTRICT;`,
	)

	if err != nil {
		logs.LogError(err, "Update Guests Team FK Query Error")

		return err
	}

	_, err = pool.Exec(
		`ALTER TABLE invites DROP CONSTRAINT invites_invitee_fkey,
		 ADD CONSTRAINT invites_invitee_fkey FOREIGN KEY(invitee)
		 REFERENCES guests(email) ON UPDATE CASCADE ON DELETE RESTRICT,
		 DROP CONSTRAINT invites_inviter_fkey,
		 ADD CONSTRAINT invites_inviter_fkey FOREIGN KEY(inviter)
		 REFERENCES admins(email) ON UPDATE CASCADE ON DELETE RESTRICT;`,
	)

	if err != nil {
		logs.LogError(err, "Update Guests Team FK Query Error")

		return err
	}

	// Update length of team id to match that generated of the unique ids created by the app.
	_, err = pool.Exec(`ALTER TABLE teams ALTER COLUMN id TYPE VARCHAR(20);`)

	if err != nil {
		logs.LogError(err, "Update Team ID Type Query Error")

		return err
	}

	_, err = pool.Exec(`ALTER TABLE admins ALTER COLUMN team TYPE VARCHAR(20);`)

	if err != nil {
		logs.LogError(err, "Update Admins Team ID Type Query Error")

		return err
	}

	_, err = pool.Exec(
		`ALTER TABLE guests ALTER COLUMN team TYPE VARCHAR(20),
		 ALTER COLUMN salt TYPE VARCHAR(10);`,
	)

	if err != nil {
		logs.LogError(err, "Update Guest Team ID Type Query Error")

		return err
	}

	_, err = pool.Exec(`ALTER TABLE uploads ALTER COLUMN team_id TYPE VARCHAR(20);`)

	if err != nil {
		logs.LogError(err, "Update Uploads Team ID Type Query Error")

		return err
	}

	return err
}

// applyMigration20230831 creates a table to track a combined list of guest and admins,
// updates existing uploads from that table, creates user roles, allows for guest admins
// to propose invites, and modifies some constraints.
func applyMigration20230831(title string) error {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	err = createAllUsersTable(pool)

	if err != nil {
		return err
	}

	err = populateAllUsersTable(pool)

	if err != nil {
		return err
	}

	err = switchUploadsUserId(pool)

	if err != nil {
		return err
	}

	err = enableUserRoles(pool)

	if err != nil {
		return err
	}

	err = addInviteProposerSupport(pool)

	if err != nil {
		return err
	}

	err = addModifiedDates(pool)

	if err != nil {
		return err
	}

	err = updateConstraints(pool)

	if err != nil {
		return err
	}

	err = recordMigration(title)

	return err
}
