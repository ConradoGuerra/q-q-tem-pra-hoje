package recipe

import (
	"q-q-tem-pra-hoje/internal/domain/ingredient"
)

type RecipeProvider interface {
	Create(Recipe) error
	GetRecommendations(ingredients *[]ingredient.Ingredient) ([]Recommendation, error)
	FindRecipes() ([]Recipe, error)
}
