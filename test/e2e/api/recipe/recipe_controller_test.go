package e2e_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"q-q-tem-pra-hoje/internal/app"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recipe"
	"q-q-tem-pra-hoje/internal/testutil"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
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
	handler := app.NewHandler(db)

	t.Cleanup(func() {
		teardown()
		db.Close()

	})

	ts := httptest.NewServer(handler)

	defer ts.Close()

	t.Run("should create a recipe", func(t *testing.T) {

		body := `{"name":"Rice", "ingredients": [
        {"name": "Onion", "measureType":"unit","quantity":1},
        {"name": "Rice", "measureType":"mg","quantity":500},
        {"name": "Garlic", "measureType":"unit","quantity":2}
      ]}`

		resp, err := http.Post(ts.URL+"/recipe", "application/json", bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("failed to get ingredients: %v", err)
		}

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
		}
	})
}

func TestRecipeController_GetRecipes(t *testing.T) {
	db, teardown := setupDatabase(t)
	handler := app.NewHandler(db)

	t.Cleanup(func() {
		teardown()
		db.Close()

	})

	ts := httptest.NewServer(handler)

	defer ts.Close()

	t.Run("should retrieve the recipes", func(t *testing.T) {

		query := `INSERT INTO ingredients_storage(name, measure_type, quantity)
            VALUES ($1, $2, $3), ($4, $5, $6);`
		_, err := db.Exec(query, "Onion", "unit", 1, "Rice", "mg", 500)

		if err != nil {
			t.Fatal(err)
		}

		id := 1
		recipes := []recipe.Recipe{
			{Id: &id, Name: "Rice with Onion and Garlic", Ingredients: []ingredient.Ingredient{
				{Name: "Onion", MeasureType: "unit", Quantity: 1},
				{Name: "Rice", MeasureType: "mg", Quantity: 500},
				{Name: "Garlic", MeasureType: "unit", Quantity: 2},
			}},
		}

		for _, recipe := range recipes {
			var recipeId int
			err := db.QueryRow("INSERT INTO recipes (name) VALUES ($1) RETURNING id;", recipe.Name).Scan(&recipeId)
			if err != nil {
				t.Fatalf("failed to insert recipe %q: %v", recipe.Name, err)
			}

			for _, ing := range recipe.Ingredients {
				_, err := db.Exec(`
                INSERT INTO recipes_ingredients (recipe_id, name, measure_type, quantity)
                VALUES ($1, $2, $3, $4)
                ON CONFLICT (recipe_id, name) DO NOTHING;
            `, recipeId, ing.Name, ing.MeasureType, ing.Quantity)
				if err != nil {
					t.Fatalf("failed to insert ingredient %q for recipe %q: %v", ing.Name, recipe.Name, err)
				}
			}
		}

		resp, err := http.Get(ts.URL + "/recipe")
		if err != nil {
			t.Fatalf("failed to get recipes: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		var recipesFound []recipe.Recipe
		err = json.Unmarshal(body, &recipesFound)
		assert.Equal(t, recipes, recipesFound)
	})
}

func TestRecipeController_Delete(t *testing.T) {
	db, teardown := setupDatabase(t)
	handler := app.NewHandler(db)

	t.Cleanup(func() {
		teardown()
		db.Close()

	})

	ts := httptest.NewServer(handler)

	defer ts.Close()

	t.Run("should delete a recipe", func(t *testing.T) {

		query := `INSERT INTO ingredients_storage(name, measure_type, quantity)
            VALUES ($1, $2, $3), ($4, $5, $6);`
		_, err := db.Exec(query, "Onion", "unit", 1, "Rice", "mg", 500)

		if err != nil {
			t.Fatal(err)
		}

		recipes := []recipe.Recipe{
			{Name: "Rice with Onion and Garlic", Ingredients: []ingredient.Ingredient{
				{Name: "Onion", MeasureType: "unit", Quantity: 1},
				{Name: "Rice", MeasureType: "mg", Quantity: 500},
				{Name: "Garlic", MeasureType: "unit", Quantity: 2},
			}},
			{Name: "Rice with Garlic", Ingredients: []ingredient.Ingredient{
				{Name: "Rice", MeasureType: "mg", Quantity: 500},
				{Name: "Garlic", MeasureType: "unit", Quantity: 2},
			}},
			{Name: "Rice with Onion", Ingredients: []ingredient.Ingredient{
				{Name: "Onion", MeasureType: "unit", Quantity: 1},
				{Name: "Rice", MeasureType: "mg", Quantity: 500},
			}},
			{Name: "Fries", Ingredients: []ingredient.Ingredient{
				{Name: "Potato", MeasureType: "unit", Quantity: 2},
			}},
		}

		for _, recipe := range recipes {
			var recipeId int
			err := db.QueryRow("INSERT INTO recipes (name) VALUES ($1) RETURNING id;", recipe.Name).Scan(&recipeId)
			if err != nil {
				t.Fatalf("failed to insert recipe %q: %v", recipe.Name, err)
			}

			for _, ing := range recipe.Ingredients {
				_, err := db.Exec(`
                INSERT INTO recipes_ingredients (recipe_id, name, measure_type, quantity)
                VALUES ($1, $2, $3, $4)
                ON CONFLICT (recipe_id, name) DO NOTHING;
            `, recipeId, ing.Name, ing.MeasureType, ing.Quantity)
				if err != nil {
					t.Fatalf("failed to insert ingredient %q for recipe %q: %v", ing.Name, recipe.Name, err)
				}
			}
		}

		req, err := http.NewRequest(http.MethodDelete, ts.URL+"/recipe?id=2", bytes.NewBufferString(`{}`))
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("failed to delete a recipe: %v", err)
		}

		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

	})
}
