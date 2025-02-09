package repositories

import (
	"q-q-tem-pra-hoje/domain/ingredient"
)

// Repo
type InMemoryIngredientRepository struct {
	Ingredients []ingredient.Ingredient
}

// Implement interface
func (r *InMemoryIngredientRepository) AddIngredient(ingredient ingredient.Ingredient) error {
	r.Ingredients = append(r.Ingredients, ingredient)
	return nil
}

func (r *InMemoryIngredientRepository) FindIngredients() ([]ingredient.Ingredient, error) {
	ingredientMap := make(map[string]ingredient.Ingredient)
	for _, ingredient := range r.Ingredients {
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
