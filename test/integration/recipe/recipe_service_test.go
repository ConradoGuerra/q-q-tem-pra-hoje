package integration_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func setupDatabase(t *testing.T) *sql.DB {
	err := godotenv.Load("../../../.env")

	if err != nil {
		t.Fatalf("Error loading .env files: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres",
		connStr)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	defer db.Close()

	return db
}
