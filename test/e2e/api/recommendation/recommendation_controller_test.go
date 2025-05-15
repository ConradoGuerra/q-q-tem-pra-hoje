package e2e_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"q-q-tem-pra-hoje/internal/app"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recipe"
	"q-q-tem-pra-hoje/internal/domain/recommendation"
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

func TestRecommendationController_GetRecommendation(t *testing.T) {
	db, teardown := setupDatabase(t)
	handler := app.NewHandler(db)

	t.Cleanup(func() {

		teardown()
		db.Close()

	})

	ts := httptest.NewServer(handler)
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

	defer ts.Close()

	t.Run("should provide the recommendations", func(t *testing.T) {

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

		resp, err := http.Get(ts.URL + "/recommendation")
		if err != nil {
			t.Fatalf("Failed to get recommendations: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		expectedRecommendations := []recommendation.Recommendation{
			{Recommendation: 1, Recipe: recipes[2]},
			{Recommendation: 2, Recipe: recipes[0]},
			{Recommendation: 3, Recipe: recipes[1]},
			{Recommendation: 4, Recipe: recipes[3]},
		}

		if err != nil {
			t.Fatalf("error while encoding expectedRecommendationsJSON: %v", err)
		}
		var recommendations []recommendation.Recommendation
		err = json.NewDecoder(resp.Body).Decode(&recommendations)
		if err != nil {
			t.Fatalf("error while decoding recommendations: %v", err)
		}
		assert.Equal(t, expectedRecommendations, recommendations)

	})
}
