package controller

import (
	"encoding/json"
	"net/http"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recipe"
)

type RecipeController struct {
	RecipeProvider recipe.RecipeProvider
}

func (rc RecipeController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		rc.Add(w, r)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Invalid HTTP Method",
	})
	return
}

func (rc RecipeController) Add(w http.ResponseWriter, r *http.Request) {

	var recipeDTO struct {
		Name        string                  `json:"name"`
		Ingredients []ingredient.Ingredient `json:"ingredients"`
	}

	if err := json.NewDecoder(r.Body).Decode(&recipeDTO); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid request body",
		})
		return
	}
	if err := rc.RecipeProvider.Create(recipe.Recipe(recipeDTO)); err != nil {
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
