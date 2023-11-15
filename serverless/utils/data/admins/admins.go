package admins

import (
	"time"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/IIP-Design/commons-gateway/utils/security/jwt"
	"github.com/IIP-Design/commons-gateway/utils/types"
	"github.com/rs/xid"
)

// CheckForActiveAdmin opens a database connection and checks whether the provided
// user email exists in the `admins` table and has the `active` value set to `true`.
func CheckForActiveAdmin(adminEmail string) (types.User, bool, error) {
	var active bool
	var inviter types.User
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	query := `SELECT email, first_name, last_name, role, team, active FROM admins WHERE email = $1;`
	err = pool.QueryRow(query, adminEmail).Scan(
		&inviter.Email, &inviter.NameFirst, &inviter.NameLast, &inviter.Role, &inviter.Team, &active)

	if err != nil {
		logs.LogError(err, "Existing Admin Query Error")
		return inviter, active, err
	}

	return inviter, active, err
}

func CheckForGuestAdmin(email string) (types.User, bool, error) {
	var active bool
	var proposer types.User
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	query := `SELECT email, first_name, last_name, role, team, expiration > NOW() AS active FROM guest_auth_data WHERE email = $1 AND role='guest admin';`
	err = pool.QueryRow(query, email).Scan(
		&proposer.Email, &proposer.NameFirst, &proposer.NameLast, &proposer.Role, &proposer.Team, &active)

	if err != nil {
		logs.LogError(err, "Guest Admin Query Error")
		return proposer, active, err
	}

	return proposer, active, err
}

// CreateAdmin opens a database connection and saves a new administrative user record.
func CreateAdmin(adminData types.User) error {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	currentTime := time.Now()

	insertAdmin :=
		`INSERT INTO admins( email, first_name, last_name, role, team, active, date_created, date_modified )
		 VALUES ( $1, $2, $3, $4, $5, $6, $7, $8 );`
	_, err = pool.Exec(insertAdmin, adminData.Email, adminData.NameFirst, adminData.NameLast, adminData.Role, adminData.Team, true, currentTime, currentTime)

	if err != nil {
		logs.LogError(err, "Create Admin Query Error")
	}

	// Add the admin to the list of all users
	guid := xid.New()

	insertAllUsers := `INSERT INTO all_users( user_id, admin_id ) VALUES ( $1, $2 );`
	_, err = pool.Exec(insertAllUsers, guid, adminData.Email)

	if err != nil {
		logs.LogError(err, "Add Admin to All Users Query Error")
	}

	return err
}

// RetrieveAdmin opens a database connection and retrieves the data for an individual admin user.
func RetrieveAdmin(username string) (map[string]any, error) {
	var admin map[string]any
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	var email string
	var first_name string
	var last_name string
	var role string
	var team string
	var active string

	query := `SELECT email, first_name, last_name, role, team, active FROM admins WHERE email = $1;`
	err = pool.QueryRow(query, username).Scan(&email, &first_name, &last_name, &role, &team, &active)

	if err != nil {
		logs.LogError(err, "Get Admin Query Error")
		return admin, err
	}

	jwt, err := jwt.GenerateJWT(username, role, false)

	if err != nil {
		logs.LogError(err, "Admin token error")
		return admin, err
	}

	admin = map[string]any{
		"email":      email,
		"givenName":  first_name,
		"familyName": last_name,
		"role":       role,
		"team":       team,
		"active":     active,
		"token":      jwt,
	}

	return admin, err
}

// RetrieveAdmins opens a database connection and retrieves the full list of admin users.
func RetrieveAdmins() ([]types.AdminUser, error) {
	var admins []types.AdminUser
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	rows, err := pool.Query(`SELECT email, first_name, last_name, role, team, active FROM admins ORDER BY first_name`)

	if err != nil {
		logs.LogError(err, "Get Admins Query Error")
		return admins, err
	}

	defer rows.Close()

	for rows.Next() {
		var admin types.AdminUser
		if err := rows.Scan(&admin.Email, &admin.NameFirst, &admin.NameLast, &admin.Role, &admin.Team, &admin.Active); err != nil {
			logs.LogError(err, "Get Admins Query Error")
			return admins, err
		}

		admins = append(admins, admin)
	}

	if err = rows.Err(); err != nil {
		logs.LogError(err, "Get Admins Query Error")
		return admins, err
	}

	return admins, err
}

// UpdateAdmin opens a database connection and updates a given
// admin user with the provided information.
// TODO? - Allow for changes to user email? If so we may need
// to add an id field and set that as the primary key on an admin.
func UpdateAdmin(admin types.AdminUser) error {
	pool := data.ConnectToDB()
	defer pool.Close()

	currentTime := time.Now()

	query :=
		`UPDATE admins SET first_name = $1, last_name = $2, role = $3, team = $4,
		 active = $5, date_modified = $6 WHERE email = $7`
	_, err := pool.Exec(query, admin.NameFirst, admin.NameLast, admin.Role, admin.Team, admin.Active, currentTime, admin.Email)

	if err != nil {
		logs.LogError(err, "Update Admin Query Error")
	}

	return err
}
