package integration_test

import (
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recipe"
	"q-q-tem-pra-hoje/internal/domain/recommendation"
	"q-q-tem-pra-hoje/internal/repository/postgres"
	recommendationService "q-q-tem-pra-hoje/internal/service/recommendation"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateRecommendations(t *testing.T) {
	db := setupDatabase(t)

	t.Cleanup(func() {
		teardownDatabase(db, t)
	})

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

	// Insert recipes and their ingredients
	for _, recipe := range recipes {
		// Insert the recipe
		var recipeID int
		err := db.QueryRow("INSERT INTO recipes (name) VALUES ($1) RETURNING id;", recipe.Name).Scan(&recipeID)
		if err != nil {
			t.Fatalf("failed to insert recipe %q: %v", recipe.Name, err)
		}

		// Insert ingredients for the recipe
		for _, ing := range recipe.Ingredients {
			_, err := db.Exec(`
                INSERT INTO recipes_ingredients (recipe_id, name, measure_type, quantity)
                VALUES ($1, $2, $3, $4)
                ON CONFLICT (recipe_id, name) DO NOTHING;
            `, recipeID, ing.Name, ing.MeasureType, ing.Quantity)
			if err != nil {
				t.Fatalf("failed to insert ingredient %q for recipe %q: %v", ing.Name, recipe.Name, err)
			}
		}
	}

	recipeManager := postgres.NewRecipeManager(db)

	service := recommendationService.NewRecommendationService(recipeManager)

	t.Run("should create the recommendations", func(t *testing.T) {

		availableIngredients := []ingredient.Ingredient{
			{Name: "Onion", MeasureType: "unit", Quantity: 1},
			{Name: "Rice", MeasureType: "mg", Quantity: 500},
		}

		recommendations, err := service.GetRecommendations(&availableIngredients)
		expectedRecommendations := []recommendation.Recommendation{
			{Recommendation: 1, Recipe: recipes[2]},
			{Recommendation: 2, Recipe: recipes[0]},
			{Recommendation: 3, Recipe: recipes[1]},
			{Recommendation: 4, Recipe: recipes[3]},
		}

		if err != nil {
			t.Errorf("error creating the recommendations: %v", err)
		}
		assert.Equal(t, expectedRecommendations, recommendations)
	})
}
