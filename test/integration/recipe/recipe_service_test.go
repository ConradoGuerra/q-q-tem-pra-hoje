package integration_test

import (
	"database/sql"
	"fmt"
	"os"
	"q-q-tem-pra-hoje/domain/ingredient"
	"q-q-tem-pra-hoje/domain/recipe"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
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

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS recipes (id SERIAL PRIMARY KEY, name TEXT NOT NULL UNIQUE);")
	if err != nil {
		t.Fatalf("Failed to creat table recipes: %v", err)
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
		t.Fatalf("Failed to creat table recipes_ingredients: %v", err)
	}

	return db
}

func teardownDatabase(db *sql.DB, t *testing.T) {
	_, err := db.Exec(`DROP TABLE recipes_ingredients`)

	if err != nil {
		t.Fatalf("Failed to drop table recipes_ingredients: %v", err)
	}

	_, err = db.Exec(`DROP TABLE recipes`)

	if err != nil {
		t.Fatalf("Failed to drop table recipes: %v", err)
	}

	db.Close()
}

type PostgreSQLRecipeManager struct {
	*sql.DB
}

func (m PostgreSQLRecipeManager) AddRecipe(recipe recipe.Recipe) error {
	var recipeID int
	err := m.QueryRow("INSERT INTO recipes (name) VALUES ($1) RETURNING id;", recipe.Name).Scan(&recipeID)
	if err != nil {
		return fmt.Errorf("Failed to insert recipe: %v", err)
	}

	for _, ing := range recipe.Ingredients {
		_, err = m.Exec(`
		      INSERT INTO recipes_ingredients (recipe_id, name, measure_type, quantity)
		      VALUES ($1, $2, $3, $4)
		      ON CONFLICT (recipe_id, name) DO NOTHING;
		  `, recipeID, ing.Name, ing.MeasureType, ing.Quantity)
		if err != nil {
			return fmt.Errorf("Failed to insert ingredient: %v", err)
		}
	}
	return nil
}
func (m PostgreSQLRecipeManager) GetAllRecipes() []recipe.Recipe {
	return []recipe.Recipe{}
}

func TestAddRecipe(t *testing.T) {
	db := setupDatabase(t)

	t.Cleanup(func() {
		teardownDatabase(db, t)
	})

	recipeManager := PostgreSQLRecipeManager{db}

	service := recipe.NewRecipeService(recipeManager)
	t.Run("should add a recipe in database", func(t *testing.T) {
		ingredients := []ingredient.Ingredient{
			{Name: "Onion", MeasureType: "unit", Quantity: 1},
			{Name: "Rice", MeasureType: "mg", Quantity: 500},
			{Name: "Garlic", MeasureType: "unit", Quantity: 2},
		}

		newRecipe := recipe.Recipe{Name: "Rice with Onion and Garlic", Ingredients: ingredients}
		err := service.AddRecipe(newRecipe)
		if err != nil {
			t.Errorf("Error on addRecipe: %v", err)
		}
		rows, err := db.Query(`SELECT 
                              r.name, 
                              i.name, 
                              i.measure_type, 
                              i.quantity 
                            FROM recipes r 
                              JOIN recipes_ingredients i ON r.id = i.recipe_id 
                            WHERE r.name = $1`, newRecipe.Name)

		if err != nil {
			t.Errorf("Error on query recipe: %v", err)
		}

		var retrievedRecipe recipe.Recipe
		for rows.Next() {
			var ing ingredient.Ingredient
			err := rows.Scan(&retrievedRecipe.Name, &ing.Name, &ing.MeasureType, &ing.Quantity)
			if err != nil {
				t.Fatalf("Failed to scan row: %v", err)
			}
			retrievedRecipe.Ingredients = append(retrievedRecipe.Ingredients, ing)
		}
		defer rows.Close()

		assert.Equal(t, newRecipe, retrievedRecipe)
	})
}
