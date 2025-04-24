package ingredient

type IngredientStorageManager interface {
	AddIngredient(Ingredient) error
	FindIngredients() ([]Ingredient, error)
	Update(Ingredient) error
	Delete(id uint) error
}
