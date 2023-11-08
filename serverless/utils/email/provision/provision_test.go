package provision

import (
	"os"
	"testing"

	testConfig "github.com/IIP-Design/commons-gateway/test/config"
	"github.com/IIP-Design/commons-gateway/utils/types"
)

func TestFormatEmail(t *testing.T) {
	testConfig.ConfigureEmail()

	invitee := types.User{
		Email:     "test@test.com",
		NameFirst: "John",
		NameLast:  "Public",
		Role:      "guest",
		Team:      "Fox",
	}

	tmpPassword := "abcfef"
	redirectUrl := os.Getenv("EMAIL_REDIRECT_URL")
	sourceEmail := os.Getenv("SOURCE_EMAIL_ADDRESS")

	e := formatEmail(
		invitee,
		tmpPassword,
		redirectUrl,
		sourceEmail,
	)

	if len(e.Destination.ToAddresses) != 1 {
		t.Fatalf(`ToAddresses length %d, want 1`, len(e.Destination.ToAddresses))
	}
	if e.Destination.ToAddresses[0] != invitee.Email {
		t.Fatalf(`ToAddresses %s, want %s`, e.Destination.ToAddresses[0], invitee.Email)
	}
}
