package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestSetup(t *testing.T) {
	err := godotenv.Load("../../../.env")

	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Connection with the database
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		fmt.Printf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            email TEXT NOT NULL UNIQUE
        )
    `)

	if err != nil {
		t.Fatalf("Failed to create the users table: %v", err)
	}

	defer func() {
		db.Exec("DROP TABLE IF EXISTS users")
	}()

}
