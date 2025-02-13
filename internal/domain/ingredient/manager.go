package ingredient

type IngredientStorageManager interface {
	AddIngredient(Ingredient) error
	FindIngredients() ([]Ingredient, error)
}
