package ingredient

type IngredientInventoryService struct {
	ingredientRepository IngredientManager
}

// Instance
func NewService(ingredientRepository IngredientManager) *IngredientInventoryService {
	return &IngredientInventoryService{ingredientRepository}
}

// Implements method
func (i *IngredientInventoryService) AddIngredientToInventory(ingredient Ingredient) {
	i.ingredientRepository.AddIngredient(ingredient)
}

func (i *IngredientInventoryService) FindIngredients() []Ingredient {
	return i.ingredientRepository.FindIngredients()
}
