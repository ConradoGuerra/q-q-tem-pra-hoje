package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

		assert.Equal(t, "Invalid request body", response["message"])

	})

	t.Run("should return unexpected error if any other error happens", func(t *testing.T) {
		mockService := MockIngredientService{CreateMock: func(i ingredient.Ingredient) error { return errors.New("Service error") }}

		controler := controller.NewIngredientController(&mockService)

		reqBody := `{"name":"Salt","measure_type":"unit","quantity":1}`
		req := httptest.NewRequest("POST", "/ingedients", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		controler.Create(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string

		json.NewDecoder(w.Body).Decode(&response)

		assert.Equal(t, "Unexpected error", response["message"])
	})

	t.Run("should return unexpected error on JSON marshaling failure", func(t *testing.T) {
		mockService := MockIngredientService{CreateMock: func(i ingredient.Ingredient) error { return fmt.Errorf("unserializable error: %v", make(chan int)) }}

		controler := controller.NewIngredientController(&mockService)

		reqBody := `{"name":"Salt","measure_type":"unit","quantity":1}`
		req := httptest.NewRequest("POST", "/ingedients", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		controler.Create(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string

		json.NewDecoder(w.Body).Decode(&response)

		assert.Equal(t, "Unexpected error", response["message"])
	})
}
