package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/server/controller/ingredient"
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
	t.Run("should create an ingredient", func(t *testing.T) {
		mockService := MockIngredientService{
			CreateMock: func(ing ingredient.Ingredient) error {
				return nil
			},
		}

		controller := controller.NewIngredientController(&mockService)
		reqBody := `{"name":"Salt","measure_type":"unit","quantity":1}`
		req := httptest.NewRequest("POST", "/ingredients", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		controller.Create(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		expectedIngredient := ingredient.Ingredient{Name: "Salt", MeasureType: "unit", Quantity: 1}
		assert.Equal(t, expectedIngredient.Name, mockService.ing.Name)
		assert.Equal(t, expectedIngredient.MeasureType, mockService.ing.MeasureType)
		assert.Equal(t, expectedIngredient.Quantity, mockService.ing.Quantity)
	})

	t.Run("should return an error if an input is invalid", func(t *testing.T) {
		mockService := MockIngredientService{
			CreateMock: func(ing ingredient.Ingredient) error {
				return nil
			},
		}

		controller := controller.NewIngredientController(&mockService)

		reqBody := `{"name":"","measure_type":"","quantity":"1"}`
		req := httptest.NewRequest("POST", "/ingredients", bytes.NewBufferString(reqBody))

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		controller.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)

		assert.NoError(t, err)

		// Verify the error message
		assert.Equal(t, "Validation failed", response["message"])

		// Verify the errors slice
		errors, ok := response["errors"].([]interface{})
		assert.True(t, ok, "Expected 'errors' to be a slice")

		// Check specific error messages
		expectedErrors := []string{
			"name cannot be empty",
			"measure_type cannot be empty",
		}
		for i, expected := range expectedErrors {
			assert.Equal(t, expected, errors[i])
		}
	})
}
