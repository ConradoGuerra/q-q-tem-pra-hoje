package ingredient

import "q-q-tem-pra-hoje/internal/domain/ingredient"

type IngredientStorageService struct {
	ingredientStorageManager ingredient.IngredientStorageManager
}

// Instance
func NewService(ingredientStorageManager ingredient.IngredientStorageManager) *IngredientStorageService {
	return &IngredientStorageService{ingredientStorageManager}
}

// Implements method
func (i *IngredientStorageService) Create(ingredient ingredient.Ingredient) error {
	return i.ingredientStorageManager.AddIngredient(ingredient)
}

func (i *IngredientStorageService) FindIngredients() ([]ingredient.Ingredient, error) {
	return i.ingredientStorageManager.FindIngredients()
}
