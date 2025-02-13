package in_memory_repository

import (
	"q-q-tem-pra-hoje/internal/domain/ingredient"
)

// Repo
type ingredientStorageManager struct {
	Ingredients []ingredient.Ingredient
}

func NewIngredientStorageManager() ingredientStorageManager{
	return ingredientStorageManager{}
}

// Implement interface
func (m *ingredientStorageManager) AddIngredient(ingredient ingredient.Ingredient) error {
	m.Ingredients = append(m.Ingredients, ingredient)
	return nil
}

func (m *ingredientStorageManager) FindIngredients() ([]ingredient.Ingredient, error) {
	ingredientMap := make(map[string]ingredient.Ingredient)
	for _, ingredient := range m.Ingredients {
		if existing, exists := ingredientMap[ingredient.Name]; exists {
			existing.Quantity += ingredient.Quantity
			ingredientMap[ingredient.Name] = existing
		} else {
			ingredientMap[ingredient.Name] = ingredient
		}
	}

	ingredientsFound := make([]ingredient.Ingredient, 0, len(ingredientMap))
	for _, ingredient := range ingredientMap {
		ingredientsFound = append(ingredientsFound, ingredient)
	}
	return ingredientsFound, nil
}
