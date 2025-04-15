package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recipe"
	"q-q-tem-pra-hoje/internal/repository/in_memory_repository"
	controller "q-q-tem-pra-hoje/internal/server/controller/recipe"
	ingredientService "q-q-tem-pra-hoje/internal/service/ingredient"
	recipeService "q-q-tem-pra-hoje/internal/service/recipe"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecipeController_Add(t *testing.T) {
	t.Run("should add a recipe", func(t *testing.T) {

		repository := in_memory_repository.NewRecipeManager([]recipe.Recipe{})
		service := recipeService.NewRecipeService(repository)
		controller := controller.RecipeController{RecipeProvider: service}

		body := `{"name":"Rice", "ingredients": [
        {"name": "Onion", "measureType":"unit","quantity":1},
        {"name": "Rice", "measureType":"mg","quantity":500},
        {"name": "Garlic", "measureType":"unit","quantity":2}
      ]}`

		r := httptest.NewRequest("POST", "/recipe", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		controller.Add(w, r)

		expectedRecipes, err := recipe.NewRecipe("Rice", []ingredient.Ingredient{
			{Name: "Onion", MeasureType: "unit", Quantity: 1},
			{Name: "Rice", MeasureType: "mg", Quantity: 500},
			{Name: "Garlic", MeasureType: "unit", Quantity: 2}})

		if err != nil {
			t.Fatalf("failed to create a recipe: %v", err)
		}

		assert.Equal(t, []recipe.Recipe{expectedRecipes}, repository.Recipes)

	})
}

func TestRecipeController_GetRecommendation(t *testing.T) {
	t.Run("should add a recipe", func(t *testing.T) {

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
		service := recipeService.NewRecipeService(recipeRepository)
		controller := controller.RecipeController{RecipeProvider: service, IngredientProvider: ingredientService}

		ingredientService.Add(ingredient.Ingredient{Name: "Onion", Quantity: 1, MeasureType: "unit"})
		ingredientService.Add(ingredient.Ingredient{Name: "Rice", Quantity: 500, MeasureType: "mg"})

		r := httptest.NewRequest("POST", "/recommendations", bytes.NewBufferString(`{}`))
		w := httptest.NewRecorder()

		controller.GetRecommendation(w, r)

		expectedRecommendations := []recipe.Recommendation{
			{Recommendation: 1, Recipe: recipes[2]},
			{Recommendation: 2, Recipe: recipes[0]},
			{Recommendation: 3, Recipe: recipes[1]},
			{Recommendation: 4, Recipe: recipes[3]},
		}

		expectedRecommendationsJSON, err := json.Marshal(expectedRecommendations)
		if err != nil {
			t.Fatalf("error while encoding expectedRecommendationsJSON: %v", err)
		}
		assert.JSONEq(t, string(expectedRecommendationsJSON), w.Body.String())

	})
}
