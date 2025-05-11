package recipe_test

import (
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recipe"
	"q-q-tem-pra-hoje/internal/repository/in_memory_repository"
	recipeService "q-q-tem-pra-hoje/internal/service/recipe"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecipeService_Add(t *testing.T) {
	t.Run("it should add a valid recipe", func(t *testing.T) {

		expectedRecipe, err := recipe.NewRecipe("Rice", []ingredient.Ingredient{
			{Name: "Onion", MeasureType: "unit", Quantity: 1},
			{Name: "Rice", MeasureType: "mg", Quantity: 500},
			{Name: "Garlic", MeasureType: "unit", Quantity: 2}})
		inMemoryRecipeManager := in_memory_repository.NewRecipeManager([]recipe.Recipe{})
		recipeService := recipeService.NewRecipeService(inMemoryRecipeManager)

		recipeService.AddRecipe(expectedRecipe)
		assert.NoError(t, err)
		assert.Equal(t, expectedRecipe, inMemoryRecipeManager.Recipes[len(inMemoryRecipeManager.Recipes)-1])
	})

	t.Run("it should return an error for an invalid name", func(t *testing.T) {
		invalidRecipe, err := recipe.NewRecipe("", []ingredient.Ingredient{ // Empty name
			{Name: "Onion", MeasureType: "unit", Quantity: 1},
			{Name: "Rice", MeasureType: "mg", Quantity: 500},
			{Name: "Garlic", MeasureType: "unit", Quantity: 2},
		})
		manager := in_memory_repository.NewRecipeManager([]recipe.Recipe{})
		service := recipeService.NewRecipeService(manager)

		service.Create(invalidRecipe)
		assert.Error(t, err)
		assert.Equal(t, "recipe name cannot be empty", err.Error())
		assert.Equal(t, recipe.Recipe{}, invalidRecipe)
	})

	t.Run("it should return an error for invalid ingredients", func(t *testing.T) {
		invalidRecipe, err := recipe.NewRecipe("Rice", []ingredient.Ingredient{}) // No ingredients
		manager := in_memory_repository.NewRecipeManager([]recipe.Recipe{})
		service := recipeService.NewRecipeService(manager)

		service.Create(invalidRecipe)
		assert.Error(t, err)
		assert.Equal(t, "recipe must have at least one ingredient", err.Error())
		assert.Equal(t, recipe.Recipe{}, invalidRecipe)

	})

}

func TestRecipesService_FindRecipes(t *testing.T) {
	t.Run("it should find all recipes", func(t *testing.T) {

		expectedRecipes := []recipe.Recipe{
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
			{Name: "Fries", Ingredients: []ingredient.Ingredient{
				{Name: "Potato", MeasureType: "unit", Quantity: 2},
			}},
		}
		repository := in_memory_repository.NewRecipeManager(expectedRecipes)
		service := recipeService.NewRecipeService(repository)

		recipes, err := service.FindRecipes()

		assert.Empty(t, err)
		assert.Equal(t, expectedRecipes, recipes)
	})
}

