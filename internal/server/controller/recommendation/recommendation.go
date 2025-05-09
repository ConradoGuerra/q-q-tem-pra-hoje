package controller

import (
	"encoding/json"
	"net/http"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recommendation"
)

type RecommendationController struct {
	IngredientProvider ingredient.IngredientStorageProvider
  RecommendationProvider recommendation.RecommendationProvider
}

func NewRecommendationController(isp ingredient.IngredientStorageProvider, rp recommendation.RecommendationProvider) *RecommendationController {
	return &RecommendationController{IngredientProvider: isp, RecommendationProvider: rp}
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

func (rc RecommendationController) GetRecommendation(w http.ResponseWriter, r *http.Request) {

	ingredients, err := rc.IngredientProvider.FindIngredients()
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"message": "No ingredients found"})
		return

	}

	recommendations, err := rc.RecommendationProvider.GetRecommendations(&ingredients)
	if err != nil {

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "No recommendations have been created"})
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&recommendations); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&recommendations)
		return
	}
}
