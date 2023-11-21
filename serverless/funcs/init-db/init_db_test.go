package main

import (
	"context"
	"testing"

	testConfig "github.com/IIP-Design/commons-gateway/test/config"
)

func TestInitDb(t *testing.T) {
	testConfig.ConfigureDb()

	resp, err := initDBHandler(context.TODO())
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("initDBHandler result %d/%v, want 200/nil", resp.StatusCode, err)
	}
}
