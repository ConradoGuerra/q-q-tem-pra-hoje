package controller_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/server/controller"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockIngredientService struct {
	ing        ingredient.Ingredient
	CreateMock func(ingredient.Ingredient) error
}

func (m *MockIngredientService) Create(ing ingredient.Ingredient) error {
	m.ing = ing
	return m.CreateMock(ing)
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
}
