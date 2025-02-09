package integration_test

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
		t.Fatalf("Error loading .env file: %v", err)
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
		t.Fatalf("Failed to connect to the database: %v", err)
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
		t.Fatalf("Failed to create the ingredients table: %v", err)
	}

	return db
}

func teardownDatabase(db *sql.DB, t *testing.T) {
	if _, err := db.Exec("DROP TABLE IF EXISTS ingredients"); err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}

	if err := db.Close(); err != nil {
		t.Fatalf("Failed to close db: %v", err)
	}
}

type PostgreSQLIngredientManager struct {
	db *sql.DB
}

func (m *PostgreSQLIngredientManager) AddIngredient(ingredient ingredient.Ingredient) error {
	query := "INSERT INTO ingredients (name, measure_type, quantity) VALUES ($1, $2, $3)"
	_, err := m.db.Exec(query, ingredient.Name, ingredient.MeasureType, ingredient.Quantity)
	if err != nil {
		return fmt.Errorf("Failed to add ingredient: %v", err)
	}
	return nil
}

func (m *PostgreSQLIngredientManager) FindIngredients() ([]ingredient.Ingredient, error) {
	query := "SELECT name, measure_type, sum(quantity) as quantity FROM ingredients GROUP BY name, measure_type;"
	rows, err := m.db.Query(query)

	if err != nil {
		return nil, fmt.Errorf("Error executing query: %v", err)
	}
	defer rows.Close()

	var ingredients []ingredient.Ingredient

	for rows.Next() {
		var ingredient ingredient.Ingredient

		err := rows.Scan(&ingredient.Name, &ingredient.MeasureType, &ingredient.Quantity)

		if err != nil {
			return nil, fmt.Errorf("Error scanning row: %v\n", err)
		}

		ingredients = append(ingredients, ingredient)
	}
	return ingredients, nil
}

func TestAddIngredientService(t *testing.T) {
	db := setupDatabase(t)

	t.Cleanup(func() {
		teardownDatabase(db, t)
	})

	ingredientManager := PostgreSQLIngredientManager{db}

	service := ingredient.NewService(&ingredientManager)

	ingredientCreated := ingredient.Ingredient{Name: "onion", Quantity: 10, MeasureType: "unit"}

	t.Run("it should add ingredients to database", func(t *testing.T) {

		err := service.AddIngredientToInventory(ingredientCreated)
		assert.NoError(t, err)

		var ingredientFound ingredient.Ingredient
		query := "SELECT name, measure_type, quantity FROM ingredients"
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

	query := `INSERT INTO ingredients (name, measure_type, quantity) 
            VALUES ($1, $2, $3), ($4, $5, $6);`
	_, err := db.Exec(query, "onion", "unit", 10, "onion", "unit", 10)

	if err != nil {
		t.Fatal(err)
	}

	t.Run("should find aggregated ingredients from the database", func(t *testing.T) {

		ingredientManager := PostgreSQLIngredientManager{db}
		ingredientService := ingredient.NewService(&ingredientManager)
		ingredientsFound, err := ingredientService.FindIngredients()

		expectedIngredients := []ingredient.Ingredient{{Name: "onion", MeasureType: "unit", Quantity: 20}}

		assert.NoError(t, err)
		assert.Equal(t, expectedIngredients, ingredientsFound)
	})
}
