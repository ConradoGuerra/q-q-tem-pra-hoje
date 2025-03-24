package controller_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"q-q-tem-pra-hoje/internal/domain/recipe"
	"q-q-tem-pra-hoje/internal/server/controller/recipe"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockedRecipeService struct {
	err func() error
}

func (mrs *MockedRecipeService) Add(rec recipe.Recipe) error {
	if err := mrs.err(); err != nil {
		return err
	}
	return nil
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
