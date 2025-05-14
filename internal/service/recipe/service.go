package recipe

import (
	"q-q-tem-pra-hoje/internal/domain/recipe"
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

func (rs *RecipeService) Delete(id uint) error {
	return rs.RecipeManager.DeleteRecipe(id)
}
