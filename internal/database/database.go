package database

import (
	"database/sql"
	"fmt"
	"q-q-tem-pra-hoje/internal/config"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	config := config.LoadDatabaseConfig()

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v\n", err)
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
		return nil, fmt.Errorf("failed to create the ingredients_storage table: %v\n", err)
	}

	return db, nil
}
