package ingredient_service

import (
	"database/sql"
	"fmt"
	"os"
	"q-q-tem-pra-hoje/domain/ingredient"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/joho/godotenv"
)

func setupDatabase(t *testing.T) *sql.DB {
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

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS ingredients (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            measure_type TEXT NOT NULL,
            quantity INT NOT NULL
        );
    `)

	if err != nil {
		t.Fatalf("Failed to create the users table: %v", err)
	}

	return db
}

type PostgreSQLIngredientManager struct {
	db *sql.DB
}

func (m *PostgreSQLIngredientManager) AddIngredient(ingredient ingredient.Ingredient) {
	query := "INSERT INTO ingredients (name, measure_type, quantity) VALUES ($1, $2, $3)"
	_, err := m.db.Exec(query, ingredient.Name, ingredient.MeasureType, ingredient.Quantity)
	if err != nil {
		panic(err)
	}
}

func (m *PostgreSQLIngredientManager) FindIngredients() []ingredient.Ingredient {
	return []ingredient.Ingredient{}
}

func TestAddIngredientService(t *testing.T) {
	db := setupDatabase(t)

	ingredientManager := PostgreSQLIngredientManager{db}

	defer db.Close()
	defer db.Exec("DROP TABLE IF EXISTS ingredients")
	service := ingredient.NewService(&ingredientManager)
	ingredientCreated := ingredient.Ingredient{Name: "onion", Quantity: 10, MeasureType: "unit"}
	service.AddIngredientToInventory(ingredientCreated)

	t.Run("it should add ingredients to database", func(t *testing.T) {
		var ingredientFound ingredient.Ingredient
		query := "SELECT name, measure_type, quantity FROM ingredients"
		err := db.QueryRow(query).Scan(&ingredientFound.Name, &ingredientFound.MeasureType, &ingredientFound.Quantity)
		if err != nil {
			t.Errorf("Failed to query ingredients: %v", err)

		}
		assert.Equal(t, ingredientCreated, ingredientFound)
	})
}
