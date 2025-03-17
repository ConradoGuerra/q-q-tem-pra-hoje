package in_memory_repository

import (
	"q-q-tem-pra-hoje/internal/domain/ingredient"
)

type ingredientStorageManager struct {
	Ingredients []ingredient.Ingredient
}

func NewIngredientStorageManager() ingredientStorageManager {
	return ingredientStorageManager{}
}

func (ism *ingredientStorageManager) AddIngredient(ingredient ingredient.Ingredient) error {
	ism.Ingredients = append(ism.Ingredients, ingredient)
	return nil
}

func (ism *ingredientStorageManager) FindIngredients() ([]ingredient.Ingredient, error) {
	ingredientMap := make(map[string]ingredient.Ingredient)
	for _, ingredient := range ism.Ingredients {
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
