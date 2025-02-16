package controller

import (
	"encoding/json"
	"net/http"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
)

type IngredientController struct {
	ingredientService ingredient.IngredientService
}

func NewIngredientController(s ingredient.IngredientService) *IngredientController {
	return &IngredientController{s}
}
func (c IngredientController) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		MeasureType string `json:"measure_type"`
		Quantity    int    `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ing := ingredient.Ingredient{
		Name:        input.Name,
		MeasureType: input.MeasureType,
		Quantity:    input.Quantity,
	}
	c.ingredientService.Create(ing)

	w.WriteHeader(http.StatusCreated)
}
