package recommendation

import (
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recipe"
	"q-q-tem-pra-hoje/internal/domain/recommendation"
	"sort"
)

type RecommendationService struct {
	recipe.RecipeManager
}

func NewRecommendationService(rm recipe.RecipeManager) *RecommendationService {
	return &RecommendationService{RecipeManager: rm}
}

func (rs *RecommendationService) GetRecommendations(ingredients *[]ingredient.Ingredient) ([]recommendation.Recommendation, error) {
	recipes, err := rs.GetAllRecipes()
	if err != nil {
		return nil, err
	}
	availableIngredientMap := make(map[string]bool)

	for _, ing := range *ingredients {
		availableIngredientMap[ing.Name] = true
	}

	type RecommendationScore struct {
		recipe recipe.Recipe
		score  float64
	}

	var scoredRecipes []RecommendationScore
	for _, recipe := range recipes {
		total := 0.0
		score := 0.0
		for _, ing := range recipe.Ingredients {
			total++
			if availableIngredientMap[ing.Name] {
				score++
			}
		}
		recommendationScore := RecommendationScore{recipe, score / total * 100}
		scoredRecipes = append(scoredRecipes, recommendationScore)
	}

	sort.Slice(scoredRecipes, func(i, j int) bool {
		return scoredRecipes[i].score > scoredRecipes[j].score
	})

	var recommendations []recommendation.Recommendation

	for i, scoredRecipe := range scoredRecipes {
		recommendations = append(recommendations, recommendation.Recommendation{Recommendation: i + 1, Recipe: scoredRecipe.recipe})
	}

	return recommendations, nil
}
