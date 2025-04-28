package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
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

func TestIngredientStorageController_GetAll(t *testing.T) {
	t.Run("should retrieve all ingredients from the repository", func(t *testing.T) {
		repository := in_memory_repository.NewIngredientStorageManager()
		ingredients := []ingredient.Ingredient{
			{Name: "Salt", MeasureType: "unit", Quantity: 1},
			{Name: "Pepper", MeasureType: "unit", Quantity: 2},
		}
		for _, ing := range ingredients {
			repository.Ingredients = append(repository.Ingredients, ing)
		}

		svc := service.NewService(&repository)
		ctrl := controller.NewIngredientController(svc)

		server := httptest.NewServer(ctrl)
		resp, err := http.Get(server.URL + "/ingredient")
		if err != nil {
			t.Fatalf("failed to retrieve ingredient: %v", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read response body: %v", err)
		}

		var responseIngredients []ingredient.Ingredient
		err = json.Unmarshal(body, &responseIngredients)
		if err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		assert.Len(t, responseIngredients, 2)
		assert.Equal(t, ingredients, responseIngredients)
	})

}
