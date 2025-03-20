package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recipe"
	"testing"

	"github.com/stretchr/testify/assert"
)

type RecipeController struct {
	recipeProvider recipe.RecipeProvider
}

func (rc RecipeController) Add(w http.ResponseWriter, r *http.Request) {

	var recipeDTO struct {
		Name        string                  `json:"name"`
		Ingredients []ingredient.Ingredient `json:"ingredients"`
	}

	json.NewDecoder(r.Body).Decode(&recipeDTO)
	rc.recipeProvider.Add(recipe.Recipe(recipeDTO))

}

type MockedRecipeService struct {
	recipes recipe.Recipe
	add     func(recipe recipe.Recipe) recipe.Recipe
}

func (mrs *MockedRecipeService) Add(rec recipe.Recipe) {
  mrs.recipes = rec

}

func TestRecipeController_Add(t *testing.T) {

	createdRecipe, _ := recipe.NewRecipe("Rice", []ingredient.Ingredient{
		{Name: "Onion", MeasureType: "unit", Quantity: 1},
	})

	t.Run("should add a recipe", func(t *testing.T) {
		service := MockedRecipeService{add: func(recipe recipe.Recipe) recipe.Recipe {
			return recipe
		}}

		controller := RecipeController{recipeProvider: &service}
		w := httptest.NewRecorder()

		requestBody := `{"name":"Rice", "ingredients": [{"name": "Onion", "measureType":"unit","quantity":1}]}`

		r := httptest.NewRequest("POST", "/recipe", bytes.NewBufferString(requestBody))
		controller.Add(w, r)

		assert.Equal(t, createdRecipe, service.recipes)
	})
}
