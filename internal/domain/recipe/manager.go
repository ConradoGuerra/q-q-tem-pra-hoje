package recipe

type RecipeManager interface {
	AddRecipe(recipe Recipe) error
  GetAllRecipes() ([]Recipe, error)
  DeleteRecipe(id uint) error
}
