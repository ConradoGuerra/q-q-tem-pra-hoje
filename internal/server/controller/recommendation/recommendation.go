package controller

import (
	"encoding/json"
	"net/http"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recipe"
)

type RecommendationController struct {
	IngredientProvider ingredient.IngredientStorageProvider
	RecipeProvider     recipe.RecipeProvider
}

func NewRecommendationController(isp ingredient.IngredientStorageProvider, rp recipe.RecipeProvider) *RecommendationController {
	return &RecommendationController{IngredientProvider: isp, RecipeProvider: rp}
}

func (rc RecommendationController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		rc.GetRecommendation(w, r)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Invalid HTTP Method",
	})
	return
}

func (rc RecommendationController) Add(w http.ResponseWriter, r *http.Request) {

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

func (rc RecommendationController) GetRecommendation(w http.ResponseWriter, r *http.Request) {

	ingredients, err := rc.IngredientProvider.FindIngredients()
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"message": "No ingredients found"})
		return

	}

	recipes, err := rc.RecipeProvider.GetRecommendations(&ingredients)
	if err != nil {

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "No recommendations have been created"})
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&recipes); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&recipes)
		return
	}
}
