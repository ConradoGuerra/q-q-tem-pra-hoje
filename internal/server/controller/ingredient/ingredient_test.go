package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	controller "q-q-tem-pra-hoje/internal/server/controller/ingredient"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockIngredientService struct {
	ing     ingredient.Ingredient
	AddMock func(ingredient.Ingredient) error
	Find    func() ([]ingredient.Ingredient, error)
}

func (mis *MockIngredientService) Add(ing ingredient.Ingredient) error {
	mis.ing = ing
	return mis.AddMock(ing)
}

func (mis *MockIngredientService) FindIngredients() ([]ingredient.Ingredient, error) {
	return mis.Find()
}

func TestIngredientController_Add(t *testing.T) {
	tests := []struct {
		method         string
		name           string
		requestBody    string
		expectedStatus int
		expectedBody   string
		mockService    MockIngredientService
	}{
		{
			method:         "POST",
			name:           "valid ingredient",
			requestBody:    `{"name":"Salt","measure_type":"unit","quantity":1}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   "",
			mockService: MockIngredientService{
				AddMock: func(ing ingredient.Ingredient) error {
					return nil
				},
			},
		},
		{
			method:         "POST",
			name:           "invalid input",
			requestBody:    `{"name":"","measure_type":"","quantity":"1"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Invalid request body"}`,
			mockService: MockIngredientService{
				AddMock: func(ing ingredient.Ingredient) error {
					return nil
				},
			},
		},
		{
			method:         "POST",
			name:           "service error",
			requestBody:    `{"name":"Salt","measure_type":"unit","quantity":1}`,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"Unexpected error"}`,
			mockService: MockIngredientService{
				AddMock: func(ing ingredient.Ingredient) error {
					return errors.New("Service error")
				},
			},
		},
		{
			method:         "GET",
			name:           "method http invalid",
			requestBody:    `{"name":"Salt","measure_type":"unit","quantity":1}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Method not allowed"}`,
			mockService: MockIngredientService{
				AddMock: func(ing ingredient.Ingredient) error {
					return errors.New("Service error")
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := controller.NewIngredientController(&tt.mockService)
			req := httptest.NewRequest(tt.method, "/ingredients", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			controller.Add(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		})
	}
}

func TestIngredientController_FindAll(t *testing.T) {
	t.Run("should find the ingredients", func(t *testing.T) {

		expectedIngredients := []ingredient.Ingredient{
			{Name: "onion", Quantity: 20, MeasureType: "unit"},
			{Name: "garlic", Quantity: 2, MeasureType: "unit"},
		}

		service := MockIngredientService{Find: func() ([]ingredient.Ingredient, error) {
			return expectedIngredients, nil
		}}

		controller := controller.NewIngredientController(&service)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/ingredients", nil)
		req.Header.Set("Content-Type", "application/json")

		expectedIngredientsJSONData, err := json.Marshal(expectedIngredients)

		if err != nil {
			t.Fatalf("fail to marshal expected ingredients: %v", err)
		}

		controller.FindAll(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, string(expectedIngredientsJSONData), w.Body.String())
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	})
}
