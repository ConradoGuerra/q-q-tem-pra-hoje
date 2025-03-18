package controller

import (
	"encoding/json"
	"net/http"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
)

type IngredientController struct {
	ingredientService ingredient.IngredientStorageProvider
}

func NewIngredientController(isp ingredient.IngredientStorageProvider) *IngredientController {
	return &IngredientController{isp}
}

func (ic IngredientController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		ic.Add(w, r)
	case http.MethodGet:
		ic.GetAll(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Method not allowed",
		})
	}
}

func (ic IngredientController) Add(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		MeasureType string `json:"measure_type"`
		Quantity    int    `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid request body",
		})
		return
	}

	ing := ingredient.NewIngredient(input.Name, input.MeasureType, input.Quantity)

	if err := ic.ingredientService.Add(ing); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Unexpected error",
		})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (ic IngredientController) GetAll(w http.ResponseWriter, r *http.Request) {

	var ingredients, _ = ic.ingredientService.FindIngredients()

	jsonIngredients, _ := json.Marshal(ingredients)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonIngredients)

}
