package propose

import (
	"os"
	"testing"

	testConfig "github.com/IIP-Design/commons-gateway/test/config"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
)

func TestFormatEmail(t *testing.T) {
	testConfig.ConfigureEmail()

	proposer := data.User{
		Email:     "test@test.com",
		NameFirst: "John",
		NameLast:  "Public",
		Role:      "guest",
		Team:      "Fox",
	}

	invitee := data.User{
		Email:     "test@test.com",
		NameFirst: "John",
		NameLast:  "Public",
		Role:      "guest",
		Team:      "Fox",
	}

	admin := data.User{
		Email:     "test@test.com",
		NameFirst: "John",
		NameLast:  "Public",
		Role:      "guest",
		Team:      "Fox",
	}

	url := os.Getenv("EMAIL_REDIRECT_URL")
	sourceEmail := os.Getenv("SOURCE_EMAIL_ADDRESS")

	e := formatEmail(
		proposer,
		invitee,
		admin,
		url,
		sourceEmail,
	)

	if len(e.Destination.ToAddresses) != 1 {
		t.Fatalf(`ToAddresses length %d, want 1`, len(e.Destination.ToAddresses))
	}
	if e.Destination.ToAddresses[0] != admin.Email {
		t.Fatalf(`ToAddresses %s, want %s`, e.Destination.ToAddresses[0], invitee.Email)
	}
}
