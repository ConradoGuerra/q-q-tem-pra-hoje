package recipe

type RecipeService struct {
	RecipeManager
}

func NewRecipeService(m RecipeManager) *RecipeService {
	return &RecipeService{RecipeManager: m}
}

func (s *RecipeService) AddRecipe(recipe Recipe) error {
	return s.RecipeManager.AddRecipe(recipe)
}
