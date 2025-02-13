package ingredient

type IngredientStorageService struct {
	ingredientStorageManager IngredientStorageManager
}

// Instance
func NewService(ingredientStorageManager IngredientStorageManager) *IngredientStorageService {
	return &IngredientStorageService{ingredientStorageManager}
}

// Implements method
func (i *IngredientStorageService) AddIngredientToStorage(ingredient Ingredient) error {
	return i.ingredientStorageManager.AddIngredient(ingredient)
}

func (i *IngredientStorageService) FindIngredients() ([]Ingredient, error) {
	return i.ingredientStorageManager.FindIngredients()
}
