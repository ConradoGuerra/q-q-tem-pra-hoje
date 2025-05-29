package e2e_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"q-q-tem-pra-hoje/internal/app"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recipe"
	"q-q-tem-pra-hoje/internal/domain/recommendation"
	"q-q-tem-pra-hoje/internal/testutil"
	"testing"
)

func intPtr(i int) *int {
	return &i
}

func TestRecommendationController_GetRecommendation(t *testing.T) {
	db := testutil.GetDB()
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	handler := app.NewHandler(db)

	ts := httptest.NewServer(handler)
	query := `INSERT INTO ingredients_storage(name, measure_type, quantity)
            VALUES ($1, $2, $3), ($4, $5, $6);`
	_, err = tx.Exec(query, "Onion", "unit", 1, "Rice", "mg", 500)

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
		err := tx.QueryRow("INSERT INTO recipes (name) VALUES ($1) RETURNING id;", recipe.Name).Scan(&recipeId)
		if err != nil {
			t.Fatalf("failed to insert recipe %q: %v", recipe.Name, err)
		}

		for _, ing := range recipe.Ingredients {
			_, err := tx.Exec(`
                INSERT INTO recipes_ingredients (recipe_id, name, measure_type, quantity)
                VALUES ($1, $2, $3, $4)
                ON CONFLICT (recipe_id, name) DO NOTHING;
            `, recipeId, ing.Name, ing.MeasureType, ing.Quantity)
			if err != nil {
				t.Fatalf("failed to insert ingredient %q for recipe %q: %v", ing.Name, recipe.Name, err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		tx.Rollback()
		ts.Close()
	})

	t.Run("should provide the recommendations", func(t *testing.T) {

		recipes := []recipe.Recipe{
			{Id: intPtr(1), Name: "Rice with Onion and Garlic", Ingredients: []ingredient.Ingredient{
				{Name: "Onion", MeasureType: "unit", Quantity: 1},
				{Name: "Rice", MeasureType: "mg", Quantity: 500},
				{Name: "Garlic", MeasureType: "unit", Quantity: 2},
			}},
			{Id: intPtr(2), Name: "Rice with Garlic", Ingredients: []ingredient.Ingredient{
				{Name: "Rice", MeasureType: "mg", Quantity: 500},
				{Name: "Garlic", MeasureType: "unit", Quantity: 2},
			}},
			{Id: intPtr(3), Name: "Rice with Onion", Ingredients: []ingredient.Ingredient{
				{Name: "Onion", MeasureType: "unit", Quantity: 1},
				{Name: "Rice", MeasureType: "mg", Quantity: 500},
			}},
			{Id: intPtr(4), Name: "Fries", Ingredients: []ingredient.Ingredient{
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
		for i, er := range expectedRecommendations {
			assert.Equal(t, er.Recommendation, recommendations[i].Recommendation)
			assert.Equal(t, er.Recipe.Name, recommendations[i].Recipe.Name)
			assert.Equal(t, er.Recipe.Ingredients, recommendations[i].Recipe.Ingredients)
		}

	})
}
