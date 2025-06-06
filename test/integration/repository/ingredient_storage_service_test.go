package respository_integration_test

import (
	"github.com/stretchr/testify/assert"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/repository/postgres"
	ingredientService "q-q-tem-pra-hoje/internal/service/ingredient"
	"q-q-tem-pra-hoje/internal/testutil"
	"testing"
)

func TestIngredientService_Add(t *testing.T) {
	db := testutil.GetDB()
	t.Cleanup(func() { cleanUpTable(t, db) })

	ingredientManager := postgres.NewIngredientStorageManager(db)
	service := ingredientService.NewService(&ingredientManager)

	ingredientCreated := ingredient.Ingredient{Name: "Salt", Quantity: 1, MeasureType: "unit"}
	secondIngredientCreated := ingredient.Ingredient{Name: "Salt", Quantity: 1, MeasureType: "unit"}

	t.Run("it should add ingredients to database", func(t *testing.T) {

		err := service.Add(ingredientCreated)
		assert.NoError(t, err)

		err = service.Add(secondIngredientCreated)
		assert.NoError(t, err)

		var ingredientFound ingredient.Ingredient
		query := "SELECT name, measure_type, quantity FROM ingredients_storage"
		err = db.QueryRow(query).Scan(&ingredientFound.Name, &ingredientFound.MeasureType, &ingredientFound.Quantity)

		assert.NoError(t, err)
		assert.Equal(t, ingredient.Ingredient{Name: "Salt", Quantity: 2, MeasureType: "unit"}, ingredientFound)
	})
}

func TestIngredientService_FindIngredients(t *testing.T) {
	db := testutil.GetDB()
	createDataset(t, db)
	t.Cleanup(func() { cleanUpTable(t, db) })

	t.Run("should find aggregated ingredients from the database", func(t *testing.T) {

		ingredientManager := postgres.NewIngredientStorageManager(db)
		ingredientService := ingredientService.NewService(&ingredientManager)
		ingredientsFound, err := ingredientService.FindIngredients()

		expectedIngredients := []ingredient.Ingredient{{Name: "onion", MeasureType: "unit", Quantity: 10}, {Name: "garlic", MeasureType: "unit", Quantity: 10}}

		assert.NoError(t, err)
		for i, expected := range expectedIngredients {
			actual := ingredientsFound[i]
			assert.Equal(t, expected.Name, actual.Name)
			assert.Equal(t, expected.MeasureType, actual.MeasureType)
			assert.Equal(t, expected.Quantity, actual.Quantity)
		}
	})
}

func TestIngredientService_Update(t *testing.T) {
	db := testutil.GetDB()
	createDataset(t, db)
	t.Cleanup(func() { cleanUpTable(t, db) })

	t.Run("should update an ingredient value", func(t *testing.T) {

		ingredientManager := postgres.NewIngredientStorageManager(db)
		ingredientService := ingredientService.NewService(&ingredientManager)

		id := 2
		updatedIngredient := ingredient.Ingredient{Id: &id, Name: "garlic", Quantity: 1, MeasureType: "unit"}

		err := ingredientService.Update(updatedIngredient)
		if err != nil {
			t.Fatalf("fail to update an ingredient %v", err)
		}

		var ingredientFound ingredient.Ingredient
		query := "SELECT name, measure_type, quantity FROM ingredients_storage WHERE id = 2"
		err = db.QueryRow(query).Scan(&ingredientFound.Name, &ingredientFound.MeasureType, &ingredientFound.Quantity)

		if err != nil {
			t.Fatalf("fail to query an ingredient %v", err)
		}

		expectedIngredient := ingredient.Ingredient{Name: "garlic", MeasureType: "unit", Quantity: 1}
		assert.Equal(t, expectedIngredient, ingredientFound)
	})

}

func TestIngredientService_Delete(t *testing.T) {
	db := testutil.GetDB()
	createDataset(t, db)
	t.Cleanup(func() { cleanUpTable(t, db) })

	t.Run("should remove an ingredient", func(t *testing.T) {

		ingredientManager := postgres.NewIngredientStorageManager(db)
		ingredientService := ingredientService.NewService(&ingredientManager)

		id := uint(2)
		err := ingredientService.Delete(id)
		if err != nil {
			t.Fatalf("fail to delete an ingredient %v", err)
		}

		var ingredientFound ingredient.Ingredient
		query := "SELECT name, measure_type, quantity FROM ingredients_storage WHERE id = 2"
		err = db.QueryRow(query).Scan(&ingredientFound.Name, &ingredientFound.MeasureType, &ingredientFound.Quantity)

		if err == nil {
			t.Fatalf("expected ingredient to be deleted, but query returned data: %v", ingredientFound)
		}
	})
}
