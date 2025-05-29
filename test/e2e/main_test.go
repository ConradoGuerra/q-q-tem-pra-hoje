package e2e_test

import (
	"os"
	"q-q-tem-pra-hoje/internal/testutil"
	"testing"
)

func TestMain(m *testing.M) {
	dsn, teardown := testutil.SetupTestDB()

	db := testutil.Connect(dsn)
	testutil.SetDB(db)
	testutil.RunMigrations(db)
	code := m.Run()

	defer teardown()

	os.Exit(code)
}
