package ingredient

type IngredientManager interface {
	AddIngredient(Ingredient) error
	FindIngredients() ([]Ingredient, error)
}
