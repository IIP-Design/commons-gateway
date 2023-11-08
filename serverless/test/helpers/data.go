package testHelpers

import (
	initDb "github.com/IIP-Design/commons-gateway/utils/init"
	"github.com/rs/xid"
)

func ExampleDbRecords() [][]string {
	return [][]string{
		{"teams", "Fox", "", "", "", "", "GPAVideo"},
		{"admins", "Fox", "admin@example.com", "John", "Public", "admin", ""},
		{"guests", "Fox", "guest@example.com", "Kristy", "Thomas", "guest", ""},
	}
}

func SetUpTestDb() (string, error) {
	var returnId string
	var err error

	teamId := xid.New()
	adminAllusersId := xid.New()
	guestAllusersId := xid.New()

	teamQuery := "INSERT INTO teams( id, team_name, aprimo_name, active, date_created, date_modified ) VALUES ($1, 'Fox', 'GPAVideo', TRUE, NOW(), NOW());"

	adminTableQuery := "INSERT INTO admins( email, first_name, last_name, role, team, active, date_created, date_modified ) VALUES ( 'admin@example.com', 'Kristy', 'Thomas', 'admin', $1, TRUE, NOW(), NOW() );"
	adminAllQuestsQuery := "INSERT INTO all_users( user_id, admin_id ) VALUES ( $1, 'admin@example.com' );"

	guestTableQuery := "INSERT INTO guests( email, first_name, last_name, role, team, date_created, date_modified ) VALUES ( 'guest@example.com', 'John', 'Public', 'guest', $1, NOW(), NOW() );"
	guestAllQuestsQuery := "INSERT INTO all_users( user_id, guest_id ) VALUES ( $1, 'guest@example.com' );"

	inviteQuery := "INSERT INTO invites( invitee, inviter, pending, date_invited, pass_hash, salt, expiration, password_reset ) VALUES ( 'guest@example.com', 'admin@example.com', FALSE, NOW(), 'abcdef', 'abcd', NOW() + INTERVAL '1 YEAR', FALSE );"

	pool := initDb.ConnectToDBInit()
	defer pool.Close()

	_, err = pool.Exec(teamQuery, teamId)
	if err != nil {
		return returnId, err
	}

	_, err = pool.Exec(adminTableQuery, teamId)
	if err != nil {
		return returnId, err
	}

	_, err = pool.Exec(adminAllQuestsQuery, adminAllusersId)
	if err != nil {
		return returnId, err
	}

	_, err = pool.Exec(guestTableQuery, teamId)
	if err != nil {
		return returnId, err
	}

	_, err = pool.Exec(guestAllQuestsQuery, guestAllusersId)
	if err != nil {
		return returnId, err
	}

	_, err = pool.Exec(inviteQuery)
	if err != nil {
		return returnId, err
	}

	returnId = teamId.String()

	return returnId, nil
}

func TearDownTestDb(teamId string) error {
	var err error

	teamQuery := "DELETE FROM teams WHERE id = $1;"

	adminTableQuery := "DELETE FROM admins WHERE email = 'admin@example.com';"
	adminAllQuestsQuery := "DELETE FROM all_users WHERE admin_id = 'admin@example.com';"

	guestTableQuery := "DELETE FROM guests WHERE email = 'guest@example.com';"
	guestAllQuestsQuery := "DELETE FROM all_users WHERE guest_id = 'guest@example.com';"

	inviteQuery := "DELETE FROM invites WHERE invitee = 'guest@example.com';"

	pool := initDb.ConnectToDBInit()
	defer pool.Close()

	_, err = pool.Exec(inviteQuery)
	if err != nil {
		return err
	}

	_, err = pool.Exec(guestAllQuestsQuery)
	if err != nil {
		return err
	}

	_, err = pool.Exec(guestTableQuery)
	if err != nil {
		return err
	}

	_, err = pool.Exec(adminAllQuestsQuery)
	if err != nil {
		return err
	}

	_, err = pool.Exec(adminTableQuery)
	if err != nil {
		return err
	}

	_, err = pool.Exec(teamQuery)
	if err != nil {
		return err
	}

	return nil
}
