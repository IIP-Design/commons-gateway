package admins

import (
	"time"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/jwt"
	"github.com/IIP-Design/commons-gateway/utils/logs"
)

// CheckForActiveAdmin opens a database connection and checks whether the provided
// user email exists in the `admins` table and has the `active` value set to `true`.
func CheckForActiveAdmin(adminEmail string) (bool, error) {
	var active bool
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	query := `SELECT active FROM admins WHERE email = $1;`
	err = pool.QueryRow(query, adminEmail).Scan(&active)

	if err != nil {
		logs.LogError(err, "Existing Admin Query Error")
		return active, err
	}

	return active, err
}

// CreateAdmin opens a database connection and saves a new administrative user record.
func CreateAdmin(adminData data.User) error {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	currentTime := time.Now()

	insertAdmin :=
		`INSERT INTO admins( email, first_name, last_name, role, team, active, date_created, date_modified )
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`
	_, err = pool.Exec(insertAdmin, adminData.Email, adminData.NameFirst, adminData.NameLast, adminData.Role, adminData.Team, true, currentTime, currentTime)

	if err != nil {
		logs.LogError(err, "Create Admin Query Error")
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

	jwt, err := jwt.GenerateJWT(username, role)

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
func RetrieveAdmins() ([]map[string]any, error) {
	var admins []map[string]any
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
		var admin data.AdminUser
		if err := rows.Scan(&admin.Email, &admin.NameFirst, &admin.NameLast, &admin.Role, &admin.Team, &admin.Active); err != nil {
			logs.LogError(err, "Get Admins Query Error")
			return admins, err
		}

		adminData := map[string]any{
			"email":      admin.Email,
			"givenName":  admin.NameFirst,
			"familyName": admin.NameLast,
			"role":       admin.Role,
			"team":       admin.Team,
			"active":     admin.Active,
		}

		admins = append(admins, adminData)
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
func UpdateAdmin(admin data.AdminUser) error {
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
