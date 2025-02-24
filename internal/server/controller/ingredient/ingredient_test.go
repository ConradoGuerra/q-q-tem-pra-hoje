package controller_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	controller "q-q-tem-pra-hoje/internal/server/controller/ingredient"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockIngredientService struct {
	ing        ingredient.Ingredient
	CreateMock func(ingredient.Ingredient) error
}

func (m *MockIngredientService) AddIngredient(ing ingredient.Ingredient) error {
	m.ing = ing
	return m.CreateMock(ing)
}

func (m *MockIngredientService) FindIngredients() ([]ingredient.Ingredient, error) {
	return []ingredient.Ingredient{{}}, nil
}

func TestIngredientController_Create(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedBody   string
		mockService    MockIngredientService
	}{
		{
			name:           "valid ingredient",
			requestBody:    `{"name":"Salt","measure_type":"unit","quantity":1}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   "",
			mockService: MockIngredientService{
				CreateMock: func(ing ingredient.Ingredient) error {
					return nil
				},
			},
		},
		{
			name:           "invalid input",
			requestBody:    `{"name":"","measure_type":"","quantity":"1"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Invalid request body"}`,
			mockService: MockIngredientService{
				CreateMock: func(ing ingredient.Ingredient) error {
					return nil
				},
			},
		},
		{
			name:           "service error",
			requestBody:    `{"name":"Salt","measure_type":"unit","quantity":1}`,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"Unexpected error"}`,
			mockService: MockIngredientService{
				CreateMock: func(ing ingredient.Ingredient) error {
					return errors.New("Service error")
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := controller.NewIngredientController(&tt.mockService)
			req := httptest.NewRequest("POST", "/ingredients", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			controller.Create(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}
