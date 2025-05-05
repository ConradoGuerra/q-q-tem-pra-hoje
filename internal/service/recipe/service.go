package recipe

import (
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recipe"
	"sort"
)

type RecipeService struct {
	recipe.RecipeManager
}

func NewRecipeService(rm recipe.RecipeManager) *RecipeService {
	return &RecipeService{RecipeManager: rm}
}

func (rs *RecipeService) Create(recipe recipe.Recipe) error {
	return rs.AddRecipe(recipe)
}

func (rs *RecipeService) FindRecipes() ([]recipe.Recipe, error) {
	recipes, err := rs.GetAllRecipes()
	if err != nil {
		return nil, err
	}
	return recipes, nil
}

func (rs *RecipeService) GetRecommendations(ingredients *[]ingredient.Ingredient) ([]recipe.Recommendation, error) {
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
		recommendationScore := RecommendationScore{recipe, score / total*100}
		scoredRecipes = append(scoredRecipes, recommendationScore)
	}

	sort.Slice(scoredRecipes, func(i, j int) bool {
		return scoredRecipes[i].score > scoredRecipes[j].score
	})

	var recommendations []recipe.Recommendation

	for i, scoredRecipe := range scoredRecipes {
		recommendations = append(recommendations, recipe.Recommendation{Recommendation: i + 1, Recipe: scoredRecipe.recipe})
	}

	return recommendations, nil
}
