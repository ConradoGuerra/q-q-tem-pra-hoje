package recipe

type RecipeProvider interface {
	Create(Recipe) error
	FindRecipes() ([]Recipe, error)
}
