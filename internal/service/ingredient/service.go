package ingredient

import "q-q-tem-pra-hoje/internal/domain/ingredient"

type IngredientStorageService struct {
	ingredientStorageManager ingredient.IngredientStorageManager
}

// Instance
func NewService(ingredientStorageManager ingredient.IngredientStorageManager) *IngredientStorageService {
	return &IngredientStorageService{ingredientStorageManager}
}

func (iss *IngredientStorageService) Add(ingredient ingredient.Ingredient) error {
	return iss.ingredientStorageManager.AddIngredient(ingredient)
}

func (iss *IngredientStorageService) FindIngredients() ([]ingredient.Ingredient, error) {
	return iss.ingredientStorageManager.FindIngredients()
}

func (iss *IngredientStorageService) Update(ingredient ingredient.Ingredient)  error {
	return iss.ingredientStorageManager.Update(ingredient)
}
