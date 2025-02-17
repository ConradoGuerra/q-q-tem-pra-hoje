package controller

import (
	"encoding/json"
	"net/http"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
)

type IngredientController struct {
	ingredientService ingredient.IngredientStorageManager
}

func NewIngredientController(s ingredient.IngredientStorageManager) *IngredientController {
	return &IngredientController{s}
}
func (c IngredientController) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		MeasureType string `json:"measure_type"`
		Quantity    int    `json:"quantity"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	ing, errors := ingredient.NewIngredient(input.Name, input.MeasureType, input.Quantity)

	if errors != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorMessages := make([]string, len(errors))
		for i, err := range errors {
			errorMessages[i] = err.Error()
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Validation failed",
			"errors":  errorMessages,
		})
		return
	}
	c.ingredientService.AddIngredient(ing)

	w.WriteHeader(http.StatusCreated)
}
