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

func (m *recipeManager) AddRecipe(recipe recipe.Recipe) error {
	m.Recipes = append(m.Recipes, recipe)
	return nil
}

func (m *recipeManager) GetAllRecipes() ([]recipe.Recipe, error) {
	return m.Recipes, nil
}
