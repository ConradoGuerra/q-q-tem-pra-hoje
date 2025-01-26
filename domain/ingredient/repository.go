package ingredient

type IngredientManager interface {
	AddIngredient(Ingredient)
	FindIngredients() []Ingredient
}
