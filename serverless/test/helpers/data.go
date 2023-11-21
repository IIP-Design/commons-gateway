package testHelpers

import (
	initDb "github.com/IIP-Design/commons-gateway/utils/init"
)

const (
	GUEST_TABLE_QUERY    = "INSERT INTO guests( email, first_name, last_name, role, team, date_created, date_modified ) VALUES ( $1, $2, $3, $4, $5, NOW(), NOW() ) ON CONFLICT ON CONSTRAINT guests_pkey DO NOTHING;"
	GUEST_ALL_USER_QUERY = "INSERT INTO all_users( user_id, guest_id ) VALUES ( $1, $2 ) ON CONFLICT ON CONSTRAINT all_users_pkey DO NOTHING;"

	INVITE_QUERY         = "INSERT INTO invites( invitee, inviter, pending, date_invited, pass_hash, salt, expiration, password_reset ) VALUES ( $1, $2, FALSE, NOW(), $3, $4, NOW() + INTERVAL '1 YEAR', FALSE );"
	INVITE_PENDING_QUERY = "INSERT INTO invites( invitee, proposer, pending, date_invited, pass_hash, salt, expiration, password_reset ) VALUES ( $1, $2, TRUE, NOW(), $3, $4, NOW() + INTERVAL '1 YEAR', FALSE );"
)

var ExampleTeam = map[string]string{
	"id":          "9m4e2mr0ui3e8a21team",
	"team_name":   "Fox",
	"aprimo_name": "GPAVideo",
}

var ExampleAdmin = map[string]string{
	"user_id":    "9m4e2mr0ui3e8a2admin",
	"email":      "admin@example.com",
	"first_name": "Kristy",
	"last_name":  "Thomas",
	"role":       "admin",
}

var ExampleGuest = map[string]string{
	"user_id":    "9m4e2mr0ui3e8a2guest",
	"email":      "guest@example.com",
	"first_name": "Maryanne",
	"last_name":  "Spier",
	"role":       "guest admin",
}

var ExampleGuest2 = map[string]string{
	"user_id":    "9m4e2mr0ui3e8aguest2",
	"email":      "guest2@example.com",
	"first_name": "Imogen",
	"last_name":  "Scott",
	"role":       "guest",
}

var ExampleCreds = map[string]string{
	"salt":      "abcd",
	"pass_hash": "abcdef",
}

func ExampleDbRecords() [][]string {
	return [][]string{
		{"teams", "Fox", "", "", "", "", "GPAVideo"},
		{"admins", "Fox", "admin@example.com", "John", "Public", "admin", ""},
		{"guests", "Fox", "guest@example.com", "Kristy", "Thomas", "guest", ""},
	}
}

func SetUpTestDb() error {
	var err error

	teamQuery := "INSERT INTO teams( id, team_name, aprimo_name, active, date_created, date_modified ) VALUES ($1, $2, $3, TRUE, NOW(), NOW()) ON CONFLICT ON CONSTRAINT teams_pkey DO NOTHING;"

	adminTableQuery := "INSERT INTO admins( email, first_name, last_name, role, team, active, date_created, date_modified ) VALUES ( $1, $2, $3, $4, $5, TRUE, NOW(), NOW() ) ON CONFLICT ON CONSTRAINT admins_pkey DO NOTHING;"
	adminAllQuestsQuery := "INSERT INTO all_users( user_id, admin_id ) VALUES ( $1, $2 ) ON CONFLICT ON CONSTRAINT all_users_pkey DO NOTHING;"

	pool := initDb.ConnectToDBInit()
	defer pool.Close()

	_, err = pool.Exec(teamQuery, ExampleTeam["id"], ExampleTeam["team_name"], ExampleTeam["aprimo_name"])
	if err != nil {
		return err
	}

	_, err = pool.Exec(adminTableQuery, ExampleAdmin["email"], ExampleAdmin["first_name"], ExampleAdmin["last_name"], ExampleAdmin["role"], ExampleTeam["id"])
	if err != nil {
		return err
	}

	_, err = pool.Exec(adminAllQuestsQuery, ExampleAdmin["user_id"], ExampleAdmin["email"])
	if err != nil {
		return err
	}

	_, err = pool.Exec(GUEST_TABLE_QUERY, ExampleGuest["email"], ExampleGuest["first_name"], ExampleGuest["last_name"], ExampleGuest["role"], ExampleTeam["id"])
	if err != nil {
		return err
	}

	_, err = pool.Exec(GUEST_ALL_USER_QUERY, ExampleGuest["user_id"], ExampleGuest["email"])
	if err != nil {
		return err
	}

	_, err = pool.Exec(INVITE_QUERY, ExampleGuest["email"], ExampleAdmin["email"], ExampleCreds["pass_hash"], ExampleCreds["salt"])
	if err != nil {
		return err
	}

	return nil
}

func TearDownTestDb() error {
	var err error

	teamQuery := "DELETE FROM teams WHERE id = $1;"

	adminTableQuery := "DELETE FROM admins WHERE email = $1;"
	adminAllQuestsQuery := "DELETE FROM all_users WHERE admin_id = $1;"

	guestTableQuery := "DELETE FROM guests WHERE email = $1;"
	guestAllQuestsQuery := "DELETE FROM all_users WHERE guest_id = $1;"

	inviteQuery := "DELETE FROM invites WHERE invitee = $1;"

	pool := initDb.ConnectToDBInit()
	defer pool.Close()

	_, err = pool.Exec(inviteQuery, ExampleGuest["email"])
	if err != nil {
		return err
	}

	_, err = pool.Exec(guestAllQuestsQuery, ExampleGuest["email"])
	if err != nil {
		return err
	}

	_, err = pool.Exec(guestTableQuery, ExampleGuest["email"])
	if err != nil {
		return err
	}

	_, err = pool.Exec(adminAllQuestsQuery, ExampleAdmin["email"])
	if err != nil {
		return err
	}

	_, err = pool.Exec(adminTableQuery, ExampleAdmin["email"])
	if err != nil {
		return err
	}

	_, err = pool.Exec(teamQuery, ExampleTeam["id"])
	if err != nil {
		return err
	}

	return nil
}

func LockAccount(email string) error {
	var err error

	pool := initDb.ConnectToDBInit()
	defer pool.Close()

	query := "UPDATE guests SET locked = true WHERE email = $1"
	_, err = pool.Exec(query, email)

	return err
}

func AddPendingGuest() error {
	var err error

	pool := initDb.ConnectToDBInit()
	defer pool.Close()

	_, err = pool.Exec(GUEST_TABLE_QUERY, ExampleGuest2["email"], ExampleGuest2["first_name"], ExampleGuest2["last_name"], ExampleGuest2["role"], ExampleTeam["id"])
	if err != nil {
		return err
	}

	_, err = pool.Exec(GUEST_ALL_USER_QUERY, ExampleGuest2["user_id"], ExampleGuest2["email"])
	if err != nil {
		return err
	}

	_, err = pool.Exec(INVITE_PENDING_QUERY, ExampleGuest2["email"], ExampleGuest["email"], ExampleCreds["pass_hash"], ExampleCreds["salt"])
	if err != nil {
		return err
	}

	return nil
}
