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

		body := `{"name":"Salt","measure_type":"unit","quantity":1}`

		repo := in_memory_repository.NewIngredientStorageManager()

		srvc := service.NewService(&repo)
		ctrl := controller.NewIngredientController(srvc)

		server := httptest.NewServer(ctrl)
		defer server.Close()

		resp, err := http.Post(server.URL+"/ingredient", "application/json", bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("failed to add ingredient: %v", err)
		}

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Len(t, repo.Ingredients, 1, "repository should contain exactly one ingredient")

		expectedIngredient := ingredient.Ingredient{Name: "Salt", MeasureType: "unit", Quantity: 1}
		assert.Equal(t, expectedIngredient, repo.Ingredients[0])

		assert.Empty(t, resp.Body)
	})
}
