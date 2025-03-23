package controller

import (
	"bytes"
	"encoding/json"
	"errors"
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

	if err := json.NewDecoder(r.Body).Decode(&recipeDTO); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid request body",
		})
		return
	}
	if err := rc.recipeProvider.Add(recipe.Recipe(recipeDTO)); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Unexpected error",
		})
    return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

type MockedRecipeService struct {
	recipes recipe.Recipe
	err     func() error
}

func (mrs *MockedRecipeService) Add(rec recipe.Recipe) error {
	if err := mrs.err(); err != nil {
		return err
	}
	mrs.recipes = rec
	return nil

}

func TestRecipeController_Add(t *testing.T) {

	t.Run("should return 201 and add a recipe", func(t *testing.T) {
		createdRecipe, _ := recipe.NewRecipe("Rice", []ingredient.Ingredient{
			{Name: "Onion", MeasureType: "unit", Quantity: 1},
		})
		service := MockedRecipeService{err: func() error { return nil }}

		controller := RecipeController{recipeProvider: &service}
		w := httptest.NewRecorder()

		requestBody := `{"name":"Rice", "ingredients": [{"name": "Onion", "measureType":"unit","quantity":1}]}`

		r := httptest.NewRequest("POST", "/recipe", bytes.NewBufferString(requestBody))
		controller.Add(w, r)

		assert.Equal(t, createdRecipe, service.recipes)
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	})

	t.Run("should return 400 and message error when the input is not valid", func(t *testing.T) {
		service := MockedRecipeService{err: func() error { return nil }}

		controller := RecipeController{recipeProvider: &service}
		w := httptest.NewRecorder()

		requestBody := `{"name":, "ingredients": [{"measureType":"","quantity":1}]}`

		r := httptest.NewRequest("POST", "/recipe", bytes.NewBufferString(requestBody))
		controller.Add(w, r)

		assert.JSONEq(t, `{"message":"Invalid request body"}`, w.Body.String())
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	})

	t.Run("should return 500 and message error when unexpected error happens", func(t *testing.T) {
		service := MockedRecipeService{err: func() error {
			return errors.New("unexpected error")
		}}

		controller := RecipeController{recipeProvider: &service}
		w := httptest.NewRecorder()

		requestBody := `{"name":"Rice", "ingredients": [{"name": "Onion", "measureType":"unit","quantity":1}]}`
		r := httptest.NewRequest("POST", "/recipe", bytes.NewBufferString(requestBody))

		controller.Add(w, r)
		assert.JSONEq(t, `{"message":"Unexpected error"}`, w.Body.String())
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	})
}
