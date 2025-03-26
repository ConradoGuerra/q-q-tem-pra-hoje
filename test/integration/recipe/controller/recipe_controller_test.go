package integration_test

import (
	"bytes"
	"net/http/httptest"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recipe"
	"q-q-tem-pra-hoje/internal/repository/in_memory_repository"
	controller "q-q-tem-pra-hoje/internal/server/controller/recipe"
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
