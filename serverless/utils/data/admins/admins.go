package admins

import (
	"fmt"
	"time"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
)

// CheckForActiveAdmin opens a database connection and checks whether the provided
// user email exists in the `admins` table and has the `active` value set to `true`.
func CheckForActiveAdmin(adminEmail string) (bool, error) {
	var active bool
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	query := fmt.Sprintf(`SELECT active FROM admins WHERE email = '%s';`, adminEmail)
	err = pool.QueryRow(query).Scan(&active)

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
		`INSERT INTO "admins"("email", "first_name", "last_name", "team", "active", "date_created")
		 VALUES ($1, $2, $3, $4, $5, $6);`
	_, err = pool.Exec(insertAdmin, adminData.Email, adminData.NameFirst, adminData.NameLast, adminData.Team, true, currentTime)

	if err != nil {
		logs.LogError(err, "Create Admin Query Error")
	}

	return err
}

// RetrieveAdmins opens a database connection and retrieves the full list of admin users.
func RetrieveAdmins() ([]map[string]any, error) {
	var admins []map[string]any
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	rows, err := pool.Query(`SELECT email, first_name, last_name, team, active FROM admins`)

	if err != nil {
		logs.LogError(err, "Get Admins Query Error")
		return admins, err
	}

	defer rows.Close()

	for rows.Next() {
		var admin data.AdminUser
		if err := rows.Scan(&admin.Email, &admin.NameFirst, &admin.NameLast, &admin.Team, &admin.Active); err != nil {
			logs.LogError(err, "Get Admins Query Error")
			return admins, err
		}

		adminData := map[string]any{
			"email":      admin.Email,
			"givenName":  admin.NameFirst,
			"familyName": admin.NameLast,
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
