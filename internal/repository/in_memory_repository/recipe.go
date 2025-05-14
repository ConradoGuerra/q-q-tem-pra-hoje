package in_memory_repository

import "q-q-tem-pra-hoje/internal/domain/recipe"

type recipeManager struct {
	Recipes      []recipe.Recipe
	MethodCalled bool
}

func NewRecipeManager(recipes []recipe.Recipe) *recipeManager {
	if len(recipes) > 0 {
		return &recipeManager{Recipes: recipes, MethodCalled: false}
	}
	return &recipeManager{MethodCalled: false}
}

func (rm *recipeManager) AddRecipe(recipe recipe.Recipe) error {
	rm.MethodCalled = true
	rm.Recipes = append(rm.Recipes, recipe)
	return nil
}

func (rm *recipeManager) GetAllRecipes() ([]recipe.Recipe, error) {
	rm.MethodCalled = true
	return rm.Recipes, nil
}

func (rm *recipeManager) DeleteRecipe(id uint) error {
	rm.MethodCalled = true
  return nil
}
