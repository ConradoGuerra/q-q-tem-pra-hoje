package controller_integration_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recipe"
	"q-q-tem-pra-hoje/internal/domain/recommendation"
	"q-q-tem-pra-hoje/internal/repository/in_memory_repository"
	controller "q-q-tem-pra-hoje/internal/server/controller/recommendation"
	ingredientService "q-q-tem-pra-hoje/internal/service/ingredient"
	recommendationService "q-q-tem-pra-hoje/internal/service/recommendation"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecommendationController_GetRecommendation(t *testing.T) {
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
		ingredientRepository := in_memory_repository.NewIngredientStorageManager()
		recipeRepository := in_memory_repository.NewRecipeManager(recipes)

		ingredientService := ingredientService.NewService(&ingredientRepository)
		service := recommendationService.NewRecommendationService(recipeRepository)
		controller := controller.RecommendationController{RecommendationProvider: service, IngredientProvider: ingredientService}

		ingredientService.Add(ingredient.Ingredient{Name: "Onion", Quantity: 1, MeasureType: "unit"})
		ingredientService.Add(ingredient.Ingredient{Name: "Rice", Quantity: 500, MeasureType: "mg"})

		mux := http.NewServeMux()
		mux.HandleFunc("/recommendation", controller.GetRecommendation)
		server := httptest.NewServer(mux)

		defer server.Close()

		resp, err := http.Get(server.URL + "/recommendation")
		if err != nil {
			t.Fatalf("failed to get recommendations %v", err)
		}
		defer resp.Body.Close()

		expectedRecommendations := []recommendation.Recommendation{
			{Recommendation: 1, Recipe: recipes[2]},
			{Recommendation: 2, Recipe: recipes[0]},
			{Recommendation: 3, Recipe: recipes[1]},
			{Recommendation: 4, Recipe: recipes[3]},
		}

		body, err := io.ReadAll(resp.Body)

		if err != nil {
			t.Fatalf("failed to read response body: %v", err)
		}

		var result []recommendation.Recommendation
		err = json.Unmarshal(body, &result)
		if err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		assert.Equal(t, expectedRecommendations, result)
		assert.Equal(t, resp.StatusCode, http.StatusOK)

	})
}
