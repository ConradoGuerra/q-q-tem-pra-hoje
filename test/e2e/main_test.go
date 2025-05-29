package e2e_test

import (
	"database/sql"
	"log"
	"os"
	"q-q-tem-pra-hoje/internal/testutil"
	"testing"

	_ "github.com/lib/pq"
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

func runMigrations(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS recipes (id SERIAL PRIMARY KEY, name TEXT NOT NULL UNIQUE);")
	if err != nil {
		log.Fatalf("failed to create table recipes: %v", err)
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS recipes_ingredients (
        recipe_id INT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
        name TEXT NOT NULL,
        measure_type TEXT NOT NULL,
        quantity INT NOT NULL,
        PRIMARY KEY (recipe_id,name));
    `)

	if err != nil {
		log.Fatalf("failed to create table recipes_ingredients: %v", err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS ingredients_storage (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            measure_type TEXT NOT NULL,
            quantity INT NOT NULL
        );
    `)

	if err != nil {
		log.Fatalf("failed to create the ingredients_storage table: %v", err)
	}
}
