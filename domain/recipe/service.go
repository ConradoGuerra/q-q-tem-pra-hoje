package recipe

import (
	"q-q-tem-pra-hoje/domain/ingredient"
)

type RecipeService struct {
	RecipeManager
}

func NewRecipeService(m RecipeManager) *RecipeService {
	return &RecipeService{RecipeManager: m}
}

func (s *RecipeService) AddRecipe(recipe Recipe) error {
	return s.RecipeManager.AddRecipe(recipe)
}

func (s *RecipeService) CreateRecipeRecommendations(ingredients *[]ingredient.Ingredient) []Recommendation {
	return []Recommendation{}
}
