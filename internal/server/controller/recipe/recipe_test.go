package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recipe"
	"q-q-tem-pra-hoje/internal/server/controller/recipe"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockedRecipeService struct {
	err             func() error
	recommendations []recipe.Recipe
}

func (mrs *MockedRecipeService) Create(rec recipe.Recipe) error {
	if err := mrs.err(); err != nil {
		return err
	}
	return nil
}

func (mrs *MockedRecipeService) GetRecommendations() []recipe.Recipe { return mrs.recommendations }

func TestRecipeController_ServeHTTP(t *testing.T) {
	t.Run("should return 400 for invalid http method", func(t *testing.T) {
		service := MockedRecipeService{err: func() error { return nil }}
		controller := controller.RecipeController{RecipeProvider: &service}

		w := httptest.NewRecorder()
		body := `{"name":"Rice", "ingredients": [{"name": "Onion", "measureType":"unit","quantity":1}]}`
		r := httptest.NewRequest("GET", "/recipe", bytes.NewBufferString(body))
		controller.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"message":"Invalid HTTP Method"}`, w.Body.String())

	})
}

func TestRecipeController_Add(t *testing.T) {
	tests := []struct {
		testCase      string
		createdRecipe recipe.Recipe
		requestBody   string
		statusCode    int
		expectedBody  string
		serviceReturn error
	}{
		{
			testCase:      "should return 201 and add a recipe",
			requestBody:   `{"name":"Rice", "ingredients": [{"name": "Onion", "measureType":"unit","quantity":1}]}`,
			statusCode:    http.StatusCreated,
			expectedBody:  "",
			serviceReturn: nil,
		},
		{
			testCase:      "should return 400 and message when the input is not valid",
			requestBody:   `{"name":, "ingredients": [{"measureType":"","quantity":1}]}`,
			statusCode:    http.StatusBadRequest,
			expectedBody:  `{"message":"Invalid request body"}`,
			serviceReturn: nil,
		},
		{
			testCase:      "should return 500 and message when unexpected error happens",
			requestBody:   `{"name": "Rice", "ingredients": [{"measureType":"unit","quantity":1}]}`,
			statusCode:    http.StatusInternalServerError,
			expectedBody:  `{"message":"Unexpected error"}`,
			serviceReturn: errors.New("unexpected error"),
		},
	}
	for _, test := range tests {

		t.Run(test.testCase, func(t *testing.T) {
			service := MockedRecipeService{err: func() error { return test.serviceReturn }}

			controller := controller.RecipeController{RecipeProvider: &service}
			w := httptest.NewRecorder()

			r := httptest.NewRequest("POST", "/recipe", bytes.NewBufferString(test.requestBody))
			controller.Add(w, r)

			if w.Body.String() != "" {
				assert.JSONEq(t, test.expectedBody, w.Body.String())
			}
			assert.Equal(t, test.statusCode, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		})
	}

}

func TestRecipeController_GetRecommendation(t *testing.T) {
	t.Run("should return the recommendations", func(t *testing.T) {
		recipes := []recipe.Recipe{
			{Name: "Rice with Garlic", Ingredients: []ingredient.Ingredient{
				{Name: "Rice", MeasureType: "mg", Quantity: 500},
				{Name: "Garlic", MeasureType: "unit", Quantity: 2},
			}},
			{Name: "Rice with Onion and Garlic", Ingredients: []ingredient.Ingredient{
				{Name: "Onion", MeasureType: "unit", Quantity: 1},
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
		service := MockedRecipeService{recommendations: recipes}
		controller := controller.RecipeController{RecipeProvider: &service}

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/recommendations", bytes.NewBufferString("{}"))
		controller.GetRecommendation(w, r)

		recipeJSON, err := json.Marshal(recipes)

		if err != nil {
			t.Errorf("fail to Marshal expectedRecipe %v", err)
		}

		assert.JSONEq(t, string(recipeJSON), w.Body.String())

	})

}
