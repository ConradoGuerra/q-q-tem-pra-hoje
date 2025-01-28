package recipe

import (
	"q-q-tem-pra-hoje/domain/ingredient"
	"sort"
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
	recipes := s.GetAllRecipes()

	availableIngredientMap := make(map[string]bool)

	for _, ing := range *ingredients {
		availableIngredientMap[ing.Name] = true
	}

	type RecommendationScore struct {
		Recipe Recipe
		Score  int
	}

	var scoredRecipes []RecommendationScore
	for _, recipe := range recipes {

		score := 0
		for _, ing := range recipe.Ingredients {

			if availableIngredientMap[ing.Name] {
				score++
			} else {
				score--
			}
		}
		recommendationScore := RecommendationScore{recipe, score}
		scoredRecipes = append(scoredRecipes, recommendationScore)
	}

	sort.Slice(scoredRecipes, func(i, j int) bool {
		return scoredRecipes[i].Score > scoredRecipes[j].Score
	})

	var recommendations []Recommendation

	for i, scoredRecipe := range scoredRecipes {
		recommendations = append(recommendations, Recommendation{Recommendation: i + 1, Recipe: scoredRecipe.Recipe})
	}

	return recommendations
}
