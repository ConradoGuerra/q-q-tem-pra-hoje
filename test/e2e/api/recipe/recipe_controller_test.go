package e2e_test

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"q-q-tem-pra-hoje/internal/app"
	"q-q-tem-pra-hoje/internal/repository/postgres"
	controller "q-q-tem-pra-hoje/internal/server/controller/recipe"
	recipeService "q-q-tem-pra-hoje/internal/service/recipe"
	"q-q-tem-pra-hoje/internal/testutil"
	"testing"

	_ "github.com/lib/pq"
)

func setupDatabase(t *testing.T) (*sql.DB, func()) {
	dsn, teardown := testutil.SetupTestDB(t)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
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
	return db, teardown
}

func TestRecipeController_Add(t *testing.T) {
	db, teardown := setupDatabase(t)
	_, err := app.NewServer(db)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {

		teardown()
		db.Close()

	})

	repository := postgres.NewRecipeManager(db)
	service := recipeService.NewRecipeService(repository)
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

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
		}
	})
}
