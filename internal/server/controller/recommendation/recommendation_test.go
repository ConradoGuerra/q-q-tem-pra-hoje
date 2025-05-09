package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recommendation"
	controller "q-q-tem-pra-hoje/internal/server/controller/recommendation"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockedRecommendationService struct {
	err                func() error
	recommendations    []recommendation.Recommendation
	hasRecommendations bool
}


func (mrs *MockedRecommendationService) GetRecommendations(ingredient *[]ingredient.Ingredient) ([]recommendation.Recommendation, error) {
	if mrs.hasRecommendations != false {
		return mrs.recommendations, nil
	}
	return nil, errors.New("any error")
}

type MockerIngredientStorageService struct {
	findIngredientsCalled bool
	hasIngedients         bool
}

func (miss *MockerIngredientStorageService) Add(ingredient.Ingredient) error {
	return nil
}

func (miss *MockerIngredientStorageService) FindIngredients() ([]ingredient.Ingredient, error) {
	miss.findIngredientsCalled = true
	if miss.hasIngedients == true {

		return nil, nil
	}
	return nil, errors.New("any error")
}

func (miss *MockerIngredientStorageService) Update(ingredient.Ingredient) error {
	return nil
}

func (miss *MockerIngredientStorageService) Delete(uint) error {
	return nil
}

func TestRecommendationController_GetRecommendation(t *testing.T) {
	t.Run("should return the recommendations", func(t *testing.T) {
		recommendations := []recommendation.Recommendation{}
		recommendationService := MockedRecommendationService{recommendations: recommendations, hasRecommendations: true}
		ingredientService := MockerIngredientStorageService{findIngredientsCalled: false, hasIngedients: true}
		controller := controller.RecommendationController{RecommendationProvider: &recommendationService, IngredientProvider: &ingredientService}

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/recommendations", bytes.NewBufferString("{}"))
		controller.GetRecommendation(w, r)

		recommendationJSON, err := json.Marshal(recommendations)

		if err != nil {
			t.Errorf("fail to Marshal expectedrecommendation %v", err)
		}
		assert.True(t, ingredientService.findIngredientsCalled)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, string(recommendationJSON), w.Body.String())

	})

	t.Run("should return a message", func(t *testing.T) {
		recommendations := []recommendation.Recommendation{}
		recommendationService := MockedRecommendationService{recommendations: recommendations, hasRecommendations: false}
		ingredientService := MockerIngredientStorageService{findIngredientsCalled: false, hasIngedients: true}
		controller := controller.RecommendationController{RecommendationProvider: &recommendationService, IngredientProvider: &ingredientService}

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/recommendations", bytes.NewBufferString("{}"))
		controller.GetRecommendation(w, r)

		assert.True(t, ingredientService.findIngredientsCalled)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message": "No recommendations have been created"}`, w.Body.String())

	})

	t.Run("should return ingredients not found", func(t *testing.T) {
		recommendations := []recommendation.Recommendation{}
		recommendationService := MockedRecommendationService{recommendations: recommendations}
		ingredientService := MockerIngredientStorageService{findIngredientsCalled: false, hasIngedients: false}
		controller := controller.RecommendationController{RecommendationProvider: &recommendationService, IngredientProvider: &ingredientService}

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/recommendations", bytes.NewBufferString("{}"))
		controller.GetRecommendation(w, r)

		assert.True(t, ingredientService.findIngredientsCalled)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.JSONEq(t, `{"message": "No ingredients found"}`, w.Body.String())

	})

}
