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

func (s *RecipeService) CreateRecipe(recipe Recipe) error {
	return s.AddRecipe(recipe)
}

func (s *RecipeService) CreateRecommendations(ingredients *[]ingredient.Ingredient) ([]Recommendation, error) {
	recipes, err := s.GetAllRecipes()
	if err != nil {
		return nil, err
	}
	availableIngredientMap := make(map[string]bool)

	for _, ing := range *ingredients {
		availableIngredientMap[ing.Name] = true
	}

	type RecommendationScore struct {
		recipe Recipe
		score  int
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
		return scoredRecipes[i].score > scoredRecipes[j].score
	})

	var recommendations []Recommendation

	for i, scoredRecipe := range scoredRecipes {
		recommendations = append(recommendations, Recommendation{Recommendation: i + 1, Recipe: scoredRecipe.recipe})
	}

	return recommendations, nil
}
