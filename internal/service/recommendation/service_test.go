package recommendation_test

import (
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recipe"
	"q-q-tem-pra-hoje/internal/domain/recommendation"
	"q-q-tem-pra-hoje/internal/repository/in_memory_repository"
	service "q-q-tem-pra-hoje/internal/service/recommendation"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecommendationService_GetRecommendations(t *testing.T) {
	t.Run("it should recommend recipes based on quantity of ingredients", func(t *testing.T) {
		availableIngredients := []ingredient.Ingredient{
			{Name: "Onion", MeasureType: "unit", Quantity: 1},
			{Name: "Rice", MeasureType: "mg", Quantity: 500},
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
			{Name: "Rice", Ingredients: []ingredient.Ingredient{
				{Name: "Rice", MeasureType: "mg", Quantity: 500},
			}},
		}
		repository := in_memory_repository.NewRecipeManager(recipes)
		service := service.NewRecommendationService(repository)

		expectedRecommendations := []recommendation.Recommendation{
			{Recommendation: 1, Recipe: recipes[2]},
			{Recommendation: 2, Recipe: recipes[4]},
			{Recommendation: 3, Recipe: recipes[0]},
			{Recommendation: 4, Recipe: recipes[1]},
			{Recommendation: 5, Recipe: recipes[3]},
		}

		recommendations, err := service.GetRecommendations(&availableIngredients)

		assert.Empty(t, err)
		assert.Equal(t, expectedRecommendations, recommendations)
	})
}
