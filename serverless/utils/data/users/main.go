package users

import (
	"database/sql"
	"fmt"

	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/IIP-Design/commons-gateway/utils/logs"
)

type UserRecord struct {
	UserId  string
	AdminId sql.NullString
	GuestId sql.NullString
	Type    string
}

// CheckIfUserExists opens a database connection and checks whether the provided
// email (which is a unique value constraint in the admins and guests tables) is
// present as either an admin_id or guest_id in the all_users table. An affirmative
// check indicates that the given user has been registered as a user in the system.
func CheckForExistingUser(email string) (bool, UserRecord, error) {
	var err error
	var user UserRecord

	pool := data.ConnectToDB()
	defer pool.Close()

	query := "SELECT user_id, admin_id, guest_id FROM all_users WHERE admin_id = $1 OR guest_id = $1;"

	err = pool.QueryRow(query, email).Scan(&user.UserId, &user.AdminId, &user.GuestId)

	if err != nil {
		// Do not return an error if no results are found.
		if err == sql.ErrNoRows {
			return false, user, nil
		}

		logs.LogError(err, "Existing User Query Error")
	}

	if user.AdminId.Valid {
		user.Type = "admin"
	} else if user.GuestId.Valid {
		user.Type = "guest"
	}

	return user.Type != "", user, err
}

// RetrieveExistingUser opens a database connection and and retrieves the user
// data for a user with the provided email from the provided table.
func retrieveExistingUser(email string, table string) (data.User, error) {
	var err error

	pool := data.ConnectToDB()
	defer pool.Close()

	var user data.User

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
			return user, nil
		}

		logs.LogError(err, "Existing User Query Error")
	}

	return user, err
}

// CheckForExistingAdminUser opens a database connection and checks whether the
// provided email belongs to an existing user and whether that user is an admin.
// If the user is indeed an admin, user data will be retrieved and returned.
func CheckForExistingAdminUser(email string) (data.User, bool, error) {
	var userData data.User

	exists, user, err := CheckForExistingUser(email)

	if err != nil {
		logs.LogError(err, "Check For User Error")
		return userData, false, err
	}

	if exists && user.Type == "admin" {
		userData, err = retrieveExistingUser(email, "admins")

		if err != nil {
			logs.LogError(err, "Admin Retrieval Error")
			return userData, false, err
		}

		return userData, true, err
	} else if exists && user.Type == "guest" {
		logs.LogError(fmt.Errorf("user %s exists, but is not an admin", email), "Admin Retrieval Error")
		return userData, false, err
	}

	// Handle any outlier cases
	return userData, false, err
}

// CheckForExistingGuestUser opens a database connection and checks whether the
// provided email belongs to an existing user and whether that user is a guest.
// If the user is indeed a guest, user data will be retrieved and returned.
func CheckForExistingGuestUser(email string) (data.User, bool, error) {
	var userData data.User

	exists, user, err := CheckForExistingUser(email)

	if err != nil {
		logs.LogError(err, "Check For User Error")
		return userData, false, err
	}

	if exists && user.Type == "guest" {
		userData, err = retrieveExistingUser(email, "guests")

		if err != nil {
			logs.LogError(err, "Guest Retrieval Error")
			return userData, false, err
		}

		return userData, true, err
	} else if exists && user.Type == "admin" {
		logs.LogError(fmt.Errorf("user %s exists, but is not a guest", email), "Guest Retrieval Error")
		return userData, false, err
	}

	// Handle any outlier cases
	return userData, false, err
}
