package repositories

import "q-q-tem-pra-hoje/domain/recipe"

type InMemoryRecipeManager struct {
	Recipes []recipe.Recipe
}

func NewInMemoryRecipeManager() *InMemoryRecipeManager {
	return &InMemoryRecipeManager{Recipes: []recipe.Recipe{}}
}

func (m *InMemoryRecipeManager) AddRecipe(recipe recipe.Recipe) error {
	m.Recipes = append(m.Recipes, recipe)
	return nil
}
