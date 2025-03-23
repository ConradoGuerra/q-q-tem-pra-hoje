package recipe

type RecipeProvider interface {
	Add(Recipe) error
}
