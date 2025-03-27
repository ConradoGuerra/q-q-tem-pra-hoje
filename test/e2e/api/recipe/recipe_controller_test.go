package e2e_test

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"q-q-tem-pra-hoje/internal/repository/postgres"
	controller "q-q-tem-pra-hoje/internal/server/controller/recipe"
	service "q-q-tem-pra-hoje/internal/service/recipe"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func setupDatabase(t *testing.T) *sql.DB {
	err := godotenv.Load("../../../../.env")

	if err != nil {
		t.Fatalf("error loading .env files: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		t.Fatalf("failed to connect to the database: %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS recipes (id SERIAL PRIMARY KEY, name TEXT NOT NULL UNIQUE);")
	if err != nil {
		t.Fatalf("failed to create table recipes: %v", err)
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
		t.Fatalf("failed to create table recipes_ingredients: %v", err)
	}

	return db
}

func TestRecipeController_Add(t *testing.T) {
	db := setupDatabase(t)
	repository := postgres.NewRecipeManager(db)
	service := service.NewRecipeService(repository)
	controller := controller.RecipeController{RecipeProvider: service}
	ts := httptest.NewServer(controller)

	defer ts.Close()

	t.Run("should create a recipe", func(t *testing.T) {

		body := `{"name":"Rice", "ingredients": [
        {"name": "Onion", "measureType":"unit","quantity":1},
        {"name": "Rice", "measureType":"mg","quantity":500},
        {"name": "Garlic", "measureType":"unit","quantity":2}
      ]}`

		resp, err := http.Post(ts.URL+"/recipe", "application/json", bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("Failed to get ingredients: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
		}
	})
}
