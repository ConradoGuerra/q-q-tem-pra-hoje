package ingredient_test

import (
	"q-q-tem-pra-hoje/domain/ingredient"
	"q-q-tem-pra-hoje/infrastructure/ingredient/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddIngredient(t *testing.T) {
	t.Run("it should add ingredients to inventory", func(t *testing.T) {

		repository := &repositories.InMemoryIngredientRepository{}
		ingredientService := ingredient.NewService(repository)

		ingredient := ingredient.Ingredient{Name: "onion", Quantity: 10, MeasureType: "unit"}
		ingredientService.AddIngredientToInventory(ingredient)

		assert.Contains(t, repository.Ingredients, ingredient, "Ingredient should be added to inventory")
	})
}

func TestFindIngredients(t *testing.T) {
	t.Run("it should find all ingredients in the inventory and aggregate them", func(t *testing.T) {
		repository := &repositories.InMemoryIngredientRepository{}
		ingredientService := ingredient.NewService(repository)

		ingredientService.AddIngredientToInventory(ingredient.Ingredient{Name: "onion", Quantity: 10, MeasureType: "unit"})
		ingredientService.AddIngredientToInventory(ingredient.Ingredient{Name: "garlic", Quantity: 2, MeasureType: "unit"})
		ingredientService.AddIngredientToInventory(ingredient.Ingredient{Name: "onion", Quantity: 10, MeasureType: "unit"})

		ingredients, err := ingredientService.FindIngredients()

		expectedIngredients := []ingredient.Ingredient{
			{Name: "onion", Quantity: 20, MeasureType: "unit"},
			{Name: "garlic", Quantity: 2, MeasureType: "unit"},
		}

		assert.NoError(t, err)
		assert.Equal(t, expectedIngredients, ingredients)
	})
}
