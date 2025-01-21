package recipe

import (
	"q-q-tem-pra-hoje/domain/ingredient"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Recipe struct {
	Name        string
	Ingredients []ingredient.Ingredient
}

type RecipeManager interface {
	AddRecipe(recipe Recipe)
}

// Injection
type RecipeService struct {
	recipeRepository (RecipeManager)
}

func NewService(r RecipeManager) *RecipeService {
	return &RecipeService{r}
}

func (r RecipeService) AddRecipe(recipe Recipe) {
	r.recipeRepository.AddRecipe(recipe)
}

type InMemoryRecipeRepository struct {
	Recipes []Recipe
}

func (i *InMemoryRecipeRepository) AddRecipe(recipe Recipe) {
	i.Recipes = append(i.Recipes, recipe)
}

func TestAddRecipe(t *testing.T) {
	t.Run("it should add a recipe", func(t *testing.T) {
		recipe := Recipe{
			Name:        "Rice",
			Ingredients: []ingredient.Ingredient{{Name: "onions", MeasureType: "units", Quantity: 1}},
		}
		recipeRepository := &InMemoryRecipeRepository{}
		recipeService := NewService(recipeRepository)

		recipeService.AddRecipe(Recipe{Name: "Rice", Ingredients: []ingredient.Ingredient{{Name: "onions", MeasureType: "units", Quantity: 1}}})

		assert.Contains(t, recipeRepository.Recipes, recipe)
	})
}
