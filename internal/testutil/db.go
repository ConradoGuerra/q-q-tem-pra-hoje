package testutil

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
)

var db *sql.DB

func SetupTestDB() (string, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatal(err)
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf("postgres://testuser:testpass@%s:%d/testdb?sslmode=disable", host, port.Int())

	return dsn, func() {
		container.Terminate(ctx)
	}
}

func Connect(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	return db

}
func SetDB(newDB *sql.DB) {
	db = newDB
}

func GetDB() *sql.DB {
	if db == nil {
		panic("DB not initialized")
	}
	return db
}

func RunMigrations(db *sql.DB) {
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
