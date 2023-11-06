package data

import (
	"testing"

	"github.com/IIP-Design/commons-gateway/test"
)

func TestConnect(t *testing.T) {
	test.AddToEnv()

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
