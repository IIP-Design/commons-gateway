package init

import (
	"fmt"

	"github.com/IIP-Design/commons-gateway/utils/data/admins"
	"github.com/IIP-Design/commons-gateway/utils/data/teams"
	"github.com/IIP-Design/commons-gateway/utils/logs"
	"github.com/IIP-Design/commons-gateway/utils/types"
)

func SeedDbRecord(rec []string) error {
	switch rec[0] {
	case "admins":
		var admin types.User

		team, err := teams.GetTeamIdByName(rec[1])

		if err != nil {
			logs.LogError(err, "Admin Creation Error - Team Not Found")
			return err
		}

		admin.Team = team
		admin.Email = rec[2]
		admin.NameFirst = rec[3]
		admin.NameLast = rec[4]
		admin.Role = rec[5]

		err = admins.CreateAdmin(admin)

		if err != nil {
			logs.LogError(err, "Admin Creation Error")
			return err
		}
	case "teams":
		err := teams.CreateTeam(rec[1], rec[6])

		if err != nil {
			logs.LogError(err, "Team Creation Error")
			return err
		}

	default:
		fmt.Printf("No case for table %s", rec[0])
	}

	return nil
}

// Return true if already seeded, false otherwise
// All other errors are logged but not acted on, which is not good
func SeedForTest(records [][]string) bool {
	for _, rec := range records {
		err := SeedDbRecord(rec)
		if err != nil {
			return true
		}
	}

	return false
}
