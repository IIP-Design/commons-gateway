package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	testConfig "github.com/IIP-Design/commons-gateway/test/config"
	"github.com/IIP-Design/commons-gateway/utils/data/data"
	"github.com/rs/xid"
)

const (
	NEW_CODE = "123456"
	OLD_CODE = "567890"
)

func TestMain(m *testing.M) {
	testConfig.ConfigureDb()

	err := seedMfaTable()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	exitVal := m.Run()

	cleanMfaTable()

	os.Exit(exitVal)
}

func TestClearMfa(t *testing.T) {
	expiredMFACodeHandler()
	updated, err := checkMfaTable()
	if !updated || err != nil {
		t.Fatalf("expiredMFACodeHandler result %t/%v, want true/nil", updated, err)
	}
}

func seedMfaTable() error {
	pool := data.ConnectToDB()
	defer pool.Close()

	newId := xid.New()
	oldId := xid.New()
	currentTime := time.Now()
	insertMfa := `INSERT INTO mfa( request_id, code, date_created ) VALUES ( $1, $2, $3 );`

	_, err := pool.Exec(insertMfa, newId.String(), NEW_CODE, currentTime)
	if err != nil {
		return err
	}

	_, err = pool.Exec(insertMfa, oldId.String(), OLD_CODE, currentTime.Add(time.Duration(-30)*time.Minute))
	return err
}

func checkMfaTable() (bool, error) {
	updated := false

	pool := data.ConnectToDB()
	defer pool.Close()

	query := `SELECT code FROM mfa WHERE code IN ( $1, $2 );`

	rows, err := pool.Query(query, NEW_CODE, OLD_CODE)
	if err != nil {
		return updated, err
	}

	count := 0
	for rows.Next() {
		count += 1
	}
	updated = (count == 1)

	return updated, err
}

func cleanMfaTable() error {
	pool := data.ConnectToDB()
	defer pool.Close()

	query := `DELETE FROM mfa WHERE code IN ( $1, $2 );`

	_, err := pool.Exec(query, NEW_CODE, OLD_CODE)
	return err
}
