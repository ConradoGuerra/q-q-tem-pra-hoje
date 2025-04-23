package ingredient_test

import (
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/repository/in_memory_repository"
	ingredientService "q-q-tem-pra-hoje/internal/service/ingredient"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIngredientService_Add(t *testing.T) {
	t.Run("it should add ingredients to inventory", func(t *testing.T) {

		repository := in_memory_repository.NewIngredientStorageManager()
		ingredientService := ingredientService.NewService(&repository)

		ingredient := ingredient.Ingredient{Name: "onion", Quantity: 10, MeasureType: "unit"}
		ingredientService.Add(ingredient)

		assert.Contains(t, repository.Ingredients, ingredient, "Ingredient should be added to inventory")
	})
}

func TestIngredientService_FindIngredients(t *testing.T) {
	t.Run("it should find all ingredients in the storage and aggregate them", func(t *testing.T) {
		repository := in_memory_repository.NewIngredientStorageManager()
		ingredientService := ingredientService.NewService(&repository)

		ingredientService.Add(ingredient.Ingredient{Name: "onion", Quantity: 10, MeasureType: "unit"})
		ingredientService.Add(ingredient.Ingredient{Name: "garlic", Quantity: 2, MeasureType: "unit"})
		ingredientService.Add(ingredient.Ingredient{Name: "onion", Quantity: 10, MeasureType: "unit"})

		ingredients, err := ingredientService.FindIngredients()

		expectedIngredients := []ingredient.Ingredient{
			{Name: "onion", Quantity: 20, MeasureType: "unit"},
			{Name: "garlic", Quantity: 2, MeasureType: "unit"},
		}

		assert.NoError(t, err)
		assert.Equal(t, expectedIngredients, ingredients)
	})
}

func TestIngredientService_Update(t *testing.T) {
	t.Run("it should update an ingredient value", func(t *testing.T) {
		repository := in_memory_repository.NewIngredientStorageManager()
		ingredientService := ingredientService.NewService(&repository)

		ingredientService.Add(ingredient.Ingredient{Name: "onion", Quantity: 10, MeasureType: "unit"})

		err := ingredientService.Update(ingredient.Ingredient{Name: "garlic", Quantity: 1, MeasureType: "unit"})

		expectedIngredients := []ingredient.Ingredient{
			{Name: "garlic", Quantity: 1, MeasureType: "unit"},
		}

		assert.NoError(t, err)
		assert.Equal(t, expectedIngredients, repository.Ingredients)
	})
}
