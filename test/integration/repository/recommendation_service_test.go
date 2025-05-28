package integration_repository_test

import (
	"github.com/stretchr/testify/assert"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recipe"
	"q-q-tem-pra-hoje/internal/domain/recommendation"
	"q-q-tem-pra-hoje/internal/repository/postgres"
	recommendationService "q-q-tem-pra-hoje/internal/service/recommendation"
	"q-q-tem-pra-hoje/internal/testutil"
	"testing"
)

func TestCreateRecommendations(t *testing.T) {
	db := testutil.GetDB()
	createDataset(t, db)
	t.Cleanup(func() { cleanUpTable(t, db) })

	recipeManager := postgres.NewRecipeManager(db)

	service := recommendationService.NewRecommendationService(recipeManager)

	t.Run("should create the recommendations", func(t *testing.T) {

		availableIngredients := []ingredient.Ingredient{
			{Name: "Onion", MeasureType: "unit", Quantity: 1},
			{Name: "Rice", MeasureType: "mg", Quantity: 500},
		}

		recommendations, err := service.GetRecommendations(&availableIngredients)

		if err != nil {
			t.Errorf("error creating the recommendations: %v", err)
		}

		expectedRecommendations := []recommendation.Recommendation{
			{
				Recommendation: 1, Recipe: recipe.Recipe{
					Name: "Rice with Onion and Garlic",
					Ingredients: []ingredient.Ingredient{{Id: (*int)(nil),
						Name: "Onion", MeasureType: "unit", Quantity: 1}, {Id: (*int)(nil),
						Name: "Rice", MeasureType: "mg", Quantity: 500}, {Id: (*int)(nil),
						Name: "Garlic", MeasureType: "unit", Quantity: 2}},
				},
			},
			{
				Recommendation: 2, Recipe: recipe.Recipe{
					Name: "Tomato Soup",
					Ingredients: []ingredient.Ingredient{{Id: (*int)(nil),
						Name: "Tomato", MeasureType: "unit", Quantity: 4}, {Id: (*int)(nil),
						Name: "Water", MeasureType: "ml", Quantity: 500}, {Id: (*int)(nil),
						Name: "Salt", MeasureType: "mg", Quantity: 10}},
				},
			},
		}

		for i, er := range expectedRecommendations {

			assert.Equal(t, er.Recommendation, recommendations[i].Recommendation)
			assert.Equal(t, er.Recipe.Name, recommendations[i].Recipe.Name)
			assert.Equal(t, er.Recipe.Ingredients, recommendations[i].Recipe.Ingredients)
		}
	})
}
