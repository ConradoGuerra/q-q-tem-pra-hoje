package integration_test

import (
	"database/sql"
	"fmt"
	"os"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/repository/postgres"
	ingredientService "q-q-tem-pra-hoje/internal/service/ingredient"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/joho/godotenv"
)

func setupDatabase(t *testing.T) *sql.DB {
	err := godotenv.Load("../../../.env")

	if err != nil {
		t.Fatalf("error loading .env file: %v", err)
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
		t.Fatalf("failed to connect to the database: %v", err)
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
		t.Fatalf("failed to create the ingredients_storage table: %v", err)
	}

	return db
}

func teardownDatabase(db *sql.DB, t *testing.T) {
	if _, err := db.Exec("DROP TABLE IF EXISTS ingredients_storage"); err != nil {
		t.Fatalf("failed to drop table: %v", err)
	}

	if err := db.Close(); err != nil {
		t.Fatalf("failed to close db: %v", err)
	}
}

func TestAddIngredientService(t *testing.T) {
	db := setupDatabase(t)

	t.Cleanup(func() {
		teardownDatabase(db, t)
	})

	ingredientManager := postgres.NewIngredientStorageManager(db)

	service := ingredientService.NewService(&ingredientManager)

	ingredientCreated := ingredient.Ingredient{Name: "onion", Quantity: 10, MeasureType: "unit"}

	t.Run("it should add ingredients to database", func(t *testing.T) {

		err := service.Add(ingredientCreated)
		assert.NoError(t, err)

		var ingredientFound ingredient.Ingredient
		query := "SELECT name, measure_type, quantity FROM ingredients_storage"
		err = db.QueryRow(query).Scan(&ingredientFound.Name, &ingredientFound.MeasureType, &ingredientFound.Quantity)

		assert.NoError(t, err)
		assert.Equal(t, ingredientCreated, ingredientFound)
	})
}

func TestFindIngredientsService(t *testing.T) {
	db := setupDatabase(t)

	t.Cleanup(func() {
		teardownDatabase(db, t)
	})

	query := `INSERT INTO ingredients_storage(name, measure_type, quantity) 
            VALUES ($1, $2, $3), ($4, $5, $6);`
	_, err := db.Exec(query, "onion", "unit", 10, "onion", "unit", 10)

	if err != nil {
		t.Fatal(err)
	}

	t.Run("should find aggregated ingredients from the database", func(t *testing.T) {

		ingredientManager := postgres.NewIngredientStorageManager(db)
		ingredientService := ingredientService.NewService(&ingredientManager)
		ingredientsFound, err := ingredientService.FindIngredients()

		expectedIngredients := []ingredient.Ingredient{{Name: "onion", MeasureType: "unit", Quantity: 20}}

		assert.NoError(t, err)
		assert.Equal(t, expectedIngredients, ingredientsFound)
	})
}
