package data

import (
	"testing"

	testConfig "github.com/IIP-Design/commons-gateway/test/config"
)

func TestConnect(t *testing.T) {
	testConfig.ConfigureDb()

	pool := ConnectToDB()
	defer pool.Close()

	if pool == nil {
		t.Fatal("ConnectToDB failed")
	}

	err := pool.Ping()
	if err != nil {
		t.Fatalf(`ConnectToDB ping error %v, want nil`, err)
	}
}
