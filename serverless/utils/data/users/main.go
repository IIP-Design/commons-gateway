package users

import (
	"database/sql"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/IIP-Design/commons-gateway/utils/types"
)

// CheckForExistingUser opens a database connection and checks whether the provided
// email (which is a unique value constraint in the admins and guests tables) is
// present in the provided table. An affirmative check indicates that the given user
// has the access implied by their presence in the table.
func CheckForExistingUser(email string, table string) (types.User, bool, error) {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	var user types.User

	var query string

	if table == "admins" {
		query = "SELECT email, first_name, last_name, role, team FROM admins WHERE email = $1;"
	} else if table == "guests" {
		query = "SELECT email, first_name, last_name, role, team FROM guests WHERE email = $1;"
	}

	err = pool.QueryRow(query, email).Scan(&user.Email, &user.NameFirst, &user.NameLast, &user.Role, &user.Team)

	if err != nil {
		// Do not return an error if no results are found.
		if err == sql.ErrNoRows {
			return user, false, nil
		}

		logs.LogError(err, "Existing User Query Error")
	}

	return user, user.Email == email, err
}
