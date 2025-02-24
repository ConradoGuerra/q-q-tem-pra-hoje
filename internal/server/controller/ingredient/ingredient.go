package controller

import (
	"encoding/json"
	"net/http"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
)

type IngredientController struct {
	ingredientService ingredient.IngredientStorageProvider
}

func NewIngredientController(p ingredient.IngredientStorageProvider) *IngredientController {
	return &IngredientController{p}
}
func (c IngredientController) Add(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		MeasureType string `json:"measure_type"`
		Quantity    int    `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid request body",
		})
		return
	}

	ing := ingredient.NewIngredient(input.Name, input.MeasureType, input.Quantity)

	if err := c.ingredientService.Add(ing); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Unexpected error",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
}
