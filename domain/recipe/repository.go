package recipe

type RecipeManager interface {
	AddRecipe(recipe Recipe) error
}
