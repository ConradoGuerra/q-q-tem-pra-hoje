package repository_integration_test

import (
	"github.com/stretchr/testify/assert"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recipe"
	"q-q-tem-pra-hoje/internal/repository/postgres"
	recipeService "q-q-tem-pra-hoje/internal/service/recipe"
	"q-q-tem-pra-hoje/internal/testutil"
	"testing"
)

func TestRecipeService_AddRecipe(t *testing.T) {
	db := testutil.GetDB()
	t.Cleanup(func() { cleanUpTable(t, db) })

	recipeManager := postgres.NewRecipeManager(db)
	service := recipeService.NewRecipeService(recipeManager)

	t.Run("should add a recipe in database", func(t *testing.T) {
		ingredients := []ingredient.Ingredient{
			{Name: "Onion", MeasureType: "unit", Quantity: 1},
			{Name: "Rice", MeasureType: "mg", Quantity: 500},
			{Name: "Garlic", MeasureType: "unit", Quantity: 2},
		}

		newRecipe := recipe.Recipe{Name: "Rice with Onion and Garlic", Ingredients: ingredients}
		err := service.AddRecipe(newRecipe)
		if err != nil {
			t.Errorf("error at addRecipe: %v", err)
		}
		rows, err := db.Query(`SELECT 
                              r.name, 
                              i.name, 
                              i.measure_type, 
                              i.quantity 
                            FROM recipes r 
                              JOIN recipes_ingredients i ON r.id = i.recipe_id 
                            WHERE r.name = $1`, newRecipe.Name)

		if err != nil {
			t.Errorf("error on query recipe: %v", err)
		}

		var retrievedRecipe recipe.Recipe
		for rows.Next() {
			var ing ingredient.Ingredient
			err := rows.Scan(&retrievedRecipe.Name, &ing.Name, &ing.MeasureType, &ing.Quantity)
			if err != nil {
				t.Fatalf("failed to scan row: %v", err)
			}
			retrievedRecipe.Ingredients = append(retrievedRecipe.Ingredients, ing)
		}
		defer rows.Close()

		assert.Equal(t, newRecipe, retrievedRecipe)
	})

	t.Run("should throws an error if the same recipe exists", func(t *testing.T) {
		ingredients := []ingredient.Ingredient{
			{Name: "Onion", MeasureType: "unit", Quantity: 1},
			{Name: "Rice", MeasureType: "mg", Quantity: 500},
			{Name: "Garlic", MeasureType: "unit", Quantity: 2},
		}

		newRecipe := recipe.Recipe{Name: "Rice with Onion and Garlic", Ingredients: ingredients}
		err := service.AddRecipe(newRecipe)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to insert recipe: pq: duplicate key value violates unique constraint")
	})
}

func TestRecipeService_GetAllRecipes(t *testing.T) {
	db := testutil.GetDB()
	cleanUpTable(t, db)
	createDataset(t, db)
	t.Cleanup(func() { cleanUpTable(t, db) })

	recipeManager := postgres.NewRecipeManager(db)
	service := recipeService.NewRecipeService(recipeManager)

	expectedRecipes := []recipe.Recipe{
		{Name: "Rice with Onion and Garlic",
			Ingredients: []ingredient.Ingredient{
				{Name: "Onion", MeasureType: "unit", Quantity: 1},
				{Name: "Rice", MeasureType: "mg", Quantity: 500},
				{Name: "Garlic", MeasureType: "unit", Quantity: 2},
			},
		},
		{
			Name: "Tomato Soup",
			Ingredients: []ingredient.Ingredient{
				{Name: "Tomato", MeasureType: "unit", Quantity: 4},
				{Name: "Water", MeasureType: "ml", Quantity: 500},
				{Name: "Salt", MeasureType: "mg", Quantity: 10},
			},
		},
	}

	t.Run("should retrieve all recipes with their ingredients", func(t *testing.T) {
		recipes, err := service.GetAllRecipes()
		if err != nil {
			t.Fatalf("failed to get all recipes: %v", err)
		}

		assert.Equal(t, len(expectedRecipes), len(recipes))

		for i, expectedRecipe := range expectedRecipes {
			assert.Equal(t, expectedRecipe.Name, recipes[i].Name)
			assert.ElementsMatch(t, expectedRecipe.Ingredients, recipes[i].Ingredients)
		}
	})

	t.Run("should return empty slice when no recipes exist", func(t *testing.T) {
		_, err := db.Exec("DELETE FROM recipes_ingredients")
		if err != nil {
			t.Fatalf("failed to clear recipes_ingredients: %v", err)
		}
		_, err = db.Exec("DELETE FROM recipes")
		if err != nil {
			t.Fatalf("failed to clear recipes: %v", err)
		}

		recipes, err := service.GetAllRecipes()
		if err != nil {
			t.Fatalf("failed to get all recipes: %v", err)
		}

		assert.Empty(t, recipes)
	})
}
func TestRecipeService_DeleteRecipe(t *testing.T) {
	db := testutil.GetDB()
	cleanUpTable(t, db)
	createDataset(t, db)
	t.Cleanup(func() { cleanUpTable(t, db) })

	recipeManager := postgres.NewRecipeManager(db)
	service := recipeService.NewRecipeService(recipeManager)

	t.Run("should delete a recipe and its ingredients", func(t *testing.T) {
		recipes, err := service.GetAllRecipes()
		if err != nil {
			t.Fatalf("failed to get all recipes: %v", err)
		}
		assert.NotEmpty(t, recipes)

		err = service.DeleteRecipe(uint(1))
		assert.NoError(t, err)

		var count int
		err = db.QueryRow(`SELECT COUNT(*) FROM recipes_ingredients WHERE recipe_id IN (SELECT id FROM recipes WHERE id = 1)`).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count, "ingredients should have been deleted")
	})

	t.Run("should return error when recipe doesn't exist", func(t *testing.T) {
		err := service.DeleteRecipe(1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "recipe not found")
	})
}
