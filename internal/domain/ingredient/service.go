package ingredient

type IngredientService interface {
	Create(ingredient Ingredient) error
}
