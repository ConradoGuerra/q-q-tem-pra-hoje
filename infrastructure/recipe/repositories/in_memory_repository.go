package repositories

import "q-q-tem-pra-hoje/domain/recipe"

type InMemoryRecipeManager struct {
	Recipes []recipe.Recipe
}

func NewInMemoryRecipeManager(recipes []recipe.Recipe) *InMemoryRecipeManager {
	if len(recipes) > 0 {
		return &InMemoryRecipeManager{Recipes: recipes}
	}
	return &InMemoryRecipeManager{}
}

func (m *InMemoryRecipeManager) AddRecipe(recipe recipe.Recipe) error {
	m.Recipes = append(m.Recipes, recipe)
	return nil
}

func (m *InMemoryRecipeManager) GetAllRecipes() []recipe.Recipe {
	return m.Recipes
}
