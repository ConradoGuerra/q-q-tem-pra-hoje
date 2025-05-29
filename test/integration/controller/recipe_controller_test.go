package controller_integration_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recipe"
	"q-q-tem-pra-hoje/internal/repository/in_memory_repository"
	controller "q-q-tem-pra-hoje/internal/server/controller/recipe"
	recipeService "q-q-tem-pra-hoje/internal/service/recipe"
	"testing"
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

		expectedRecipes, err := recipe.NewRecipe(0, "Rice", []ingredient.Ingredient{
			{Name: "Onion", MeasureType: "unit", Quantity: 1},
			{Name: "Rice", MeasureType: "mg", Quantity: 500},
			{Name: "Garlic", MeasureType: "unit", Quantity: 2}})

		if err != nil {
			t.Fatalf("failed to create a recipe: %v", err)
		}

		assert.Equal(t, []recipe.Recipe{expectedRecipes}, repository.Recipes)
		assert.True(t, repository.MethodCalled)

	})
}

func TestRecipeController_Get(t *testing.T) {
	t.Run("should get the recipes", func(t *testing.T) {

		expectedRecipes := []recipe.Recipe{
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
		repository := in_memory_repository.NewRecipeManager(expectedRecipes)
		service := recipeService.NewRecipeService(repository)
		controller := controller.RecipeController{RecipeProvider: service}
		server := httptest.NewServer(controller)
		defer server.Close()

		resp, err := http.Get(server.URL + "/recipes")

		if err != nil {
			t.Fatalf("failed to find a recipe: %v", err)
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)

		if err != nil {
			t.Fatalf("failed to read response body: %v", err)
		}

		var resultRecipes []recipe.Recipe
		err = json.Unmarshal(body, &resultRecipes)
		assert.Equal(t, expectedRecipes, resultRecipes)
		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.True(t, repository.MethodCalled)

	})
}

func TestRecipeController_Delete(t *testing.T) {
	t.Run("should delete a recipe by Id", func(t *testing.T) {
		repository := in_memory_repository.NewRecipeManager([]recipe.Recipe{})
		service := recipeService.NewRecipeService(repository)
		ctrl := controller.RecipeController{RecipeProvider: service}
		server := httptest.NewServer(ctrl)

		req, err := http.NewRequest(http.MethodDelete, server.URL+"/recipes?id=1", bytes.NewBufferString(``))

		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		assert.True(t, repository.MethodCalled)

	})
}
