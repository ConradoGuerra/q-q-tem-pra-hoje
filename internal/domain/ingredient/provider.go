package ingredient

type IngredientStorageProvider interface {
	Add(Ingredient) error
	FindIngredients() ([]Ingredient, error)
  Update(Ingredient) error
}
