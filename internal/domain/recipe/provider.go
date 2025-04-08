package recipe

type RecipeProvider interface {
	Create(Recipe) error
	GetRecommendations() []Recipe
}
