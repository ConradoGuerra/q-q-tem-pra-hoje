package recipe_test

import (
	"q-q-tem-pra-hoje/domain/ingredient"
	"q-q-tem-pra-hoje/domain/recipe"
	"q-q-tem-pra-hoje/infrastructure/recipe/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddRecipe(t *testing.T) {
	t.Run("It should add a valid recipe", func(t *testing.T) {

		expectedRecipe, err := recipe.NewRecipe("Rice", []ingredient.Ingredient{
			{Name: "Onion", MeasureType: "unit", Quantity: 1},
			{Name: "Rice", MeasureType: "mg", Quantity: 500},
			{Name: "Garlic", MeasureType: "unit", Quantity: 2}})
		inMemoryRecipeManager := repositories.NewInMemoryRecipeManager()
		recipeService := recipe.NewRecipeService(inMemoryRecipeManager)

		recipeService.AddRecipe(expectedRecipe)
		assert.NoError(t, err)
		assert.Equal(t, expectedRecipe, inMemoryRecipeManager.Recipes[len(inMemoryRecipeManager.Recipes)-1])
	})

	t.Run("It should return an error for an invalid name", func(t *testing.T) {
		invalidRecipe, err := recipe.NewRecipe("", []ingredient.Ingredient{ // Empty name
			{Name: "Onion", MeasureType: "unit", Quantity: 1},
			{Name: "Rice", MeasureType: "mg", Quantity: 500},
			{Name: "Garlic", MeasureType: "unit", Quantity: 2},
		})
		manager := repositories.NewInMemoryRecipeManager()
		service := recipe.NewRecipeService(manager)

		service.AddRecipe(invalidRecipe)
		assert.Error(t, err)
		assert.Equal(t, "recipe name cannot be empty", err.Error())
		assert.Equal(t, recipe.Recipe{}, invalidRecipe)
	})

	t.Run("It should return an error for invalid ingredients", func(t *testing.T) {
		invalidRecipe, err := recipe.NewRecipe("Rice", []ingredient.Ingredient{}) // No ingredients
		manager := repositories.NewInMemoryRecipeManager()
		service := recipe.NewRecipeService(manager)

		service.AddRecipe(invalidRecipe)
		assert.Error(t, err)
		assert.Equal(t, "recipe must have at least one ingredient", err.Error())
		assert.Equal(t, recipe.Recipe{}, invalidRecipe)

	})

}

func TestRecommendRecipes(t *testing.T) {
	t.Run("It should recommend recipes based on quantity of ingredients", func(t *testing.T) {
		availableIngredients := []ingredient.Ingredient{
			{Name: "Onion", MeasureType: "unit", Quantity: 1},
			{Name: "Rice", MeasureType: "mg", Quantity: 500},
			{Name: "Garlic", MeasureType: "unit", Quantity: 2},
		}

		recipes := []recipe.Recipe{
			{Name: "Rice with Onion and Garlic", Ingredients: []ingredient.Ingredient{
				{Name: "Onion", MeasureType: "unit", Quantity: 1},
				{Name: "Rice", MeasureType: "mg", Quantity: 500},
				{Name: "Garlic", MeasureType: "unit", Quantity: 2},
			}},
			{Name: "Rice with Garlic", Ingredients: []ingredient.Ingredient{
				{Name: "Rice", MeasureType: "mg", Quantity: 500},
				{Name: "Garlic", MeasureType: "unit", Quantity: 2},
			}},
			{Name: "Rice with Onion", Ingredients: []ingredient.Ingredient{
				{Name: "Onion", MeasureType: "unit", Quantity: 1},
				{Name: "Rice", MeasureType: "mg", Quantity: 500},
			}},
		}
		repository := repositories.NewInMemoryRecipeManager()
		service := recipe.NewRecipeService(repository)

		expectedRecommendations := []recipe.Recommendation{
			{Recommendation: 1, Recipe: recipes[0]},
			{Recommendation: 2, Recipe: recipes[1]},
			{Recommendation: 3, Recipe: recipes[2]},
    }

		recommendations := service.CreateRecipeRecommendations(&availableIngredients)
		assert.Equal(t, expectedRecommendations, recommendations)
	})
}
