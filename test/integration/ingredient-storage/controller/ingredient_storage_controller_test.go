package integration_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/repository/in_memory_repository"
	controller "q-q-tem-pra-hoje/internal/server/controller/ingredient"
	service "q-q-tem-pra-hoje/internal/service/ingredient"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIngredientStorageController_Add(t *testing.T) {
	t.Run("should call the service properly", func(t *testing.T) {
		var repository = in_memory_repository.NewIngredientStorageManager()
		var service = service.NewService(&repository)
		var controller = controller.NewIngredientController(service)

		var w = httptest.NewRecorder()
		var req = httptest.NewRequest("POST", "/ingredients", bytes.NewBufferString(`{"name":"Salt","measure_type":"unit","quantity":1}`))
		req.Header.Set("Content-Type", "application/json")

		controller.Add(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Len(t, repository.Ingredients, 1, "repository should contain exactly one ingredient")

		expectedIngredient := ingredient.Ingredient{Name: "Salt", MeasureType: "unit", Quantity: 1}
		assert.Equal(t, expectedIngredient, repository.Ingredients[0])

		assert.Empty(t, w.Body)
	})
}
