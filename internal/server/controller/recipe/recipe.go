package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"strconv"

	"q-q-tem-pra-hoje/internal/domain/recipe"
)

var (
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrMethodNotAllowed   = errors.New("method not allowed")
	ErrInvalidId          = errors.New("invalid id parameter")
)

type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type RecipeController struct {
	IngredientProvider ingredient.IngredientStorageProvider
	RecipeProvider     recipe.RecipeProvider
}

func NewRecipeController(isp ingredient.IngredientStorageProvider, rp recipe.RecipeProvider) *RecipeController {
	return &RecipeController{IngredientProvider: isp, RecipeProvider: rp}
}

func (rc RecipeController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		rc.Add(w, r)
		return
	}
	if r.Method == "GET" {
		rc.GetRecipes(w, r)
		return
	}
	if r.Method == "DELETE" {
		rc.Delete(w, r)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Invalid HTTP Method",
	})
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
	recipeCreated, err := recipe.NewRecipe(0, recipeDTO.Name, recipeDTO.Ingredients)

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid request body",
		})
		return
	}

	if err := rc.RecipeProvider.Create(recipe.Recipe(recipeCreated)); err != nil {
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

func (rc RecipeController) GetRecipes(w http.ResponseWriter, r *http.Request) {

	recipes, err := rc.RecipeProvider.FindRecipes()
	if err != nil {

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "No recipes have been found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&recipes)
}

func (rc RecipeController) Delete(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		rc.respondWithError(w, http.StatusBadRequest, ErrInvalidId)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		rc.respondWithError(w, http.StatusBadRequest, ErrInvalidId)
		return
	}

	if err := rc.RecipeProvider.Delete(uint(id)); err != nil {
		rc.respondWithError(w, http.StatusInternalServerError, errors.New("failed to delete recipe"))
		return
	}
	rc.respondWithJSON(w, http.StatusNoContent, nil)
}

func (rc RecipeController) respondWithError(w http.ResponseWriter, code int, err error) {
	rc.respondWithJSON(w, code, Response{Message: err.Error()})
}

func (rc RecipeController) respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
