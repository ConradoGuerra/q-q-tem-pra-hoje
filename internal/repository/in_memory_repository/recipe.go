package in_memory_repository

import "q-q-tem-pra-hoje/internal/domain/recipe"

type recipeManager struct {
	Recipes []recipe.Recipe
}

func NewRecipeManager(recipes []recipe.Recipe) *recipeManager {
	if len(recipes) > 0 {
		return &recipeManager{Recipes: recipes}
	}
	return &recipeManager{}
}

func (rm *recipeManager) AddRecipe(recipe recipe.Recipe) error {
	rm.Recipes = append(rm.Recipes, recipe)
	return nil
}

func (rm *recipeManager) GetAllRecipes() ([]recipe.Recipe, error) {
	return rm.Recipes, nil
}
