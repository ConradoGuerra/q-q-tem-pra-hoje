package integration_test

import (
	"database/sql"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recipe"
	"q-q-tem-pra-hoje/internal/repository/postgres"
	recipeService "q-q-tem-pra-hoje/internal/service/recipe"
	"q-q-tem-pra-hoje/internal/testutil"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestRecipeService_AddRecipe(t *testing.T) {
	dsn, teardown := testutil.SetupTestDB(t)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS recipes (id SERIAL PRIMARY KEY, name TEXT NOT NULL UNIQUE);")
	if err != nil {
		t.Fatalf("failed to create table recipes: %v", err)
	}
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS recipes_ingredients (
        recipe_id INT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
        name TEXT NOT NULL,
        measure_type TEXT NOT NULL,
        quantity INT NOT NULL,
        PRIMARY KEY (recipe_id,name));
    `)

	if err != nil {
		t.Fatalf("failed to create table recipes_ingredients: %v", err)
	}

	defer teardown()
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
	dsn, teardown := testutil.SetupTestDB(t)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS recipes (id SERIAL PRIMARY KEY, name TEXT NOT NULL UNIQUE);")
	if err != nil {
		t.Fatalf("failed to create table recipes: %v", err)
	}
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS recipes_ingredients (
        recipe_id INT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
        name TEXT NOT NULL,
        measure_type TEXT NOT NULL,
        quantity INT NOT NULL,
        PRIMARY KEY (recipe_id,name));
    `)
	if err != nil {
		t.Fatalf("failed to create table recipes_ingredients: %v", err)
	}

	defer teardown()
	recipeManager := postgres.NewRecipeManager(db)
	service := recipeService.NewRecipeService(recipeManager)

	testRecipes := []struct {
		recipe      recipe.Recipe
		shouldExist bool
	}{
		{
			recipe: recipe.Recipe{
				Name: "Rice with Onion and Garlic",
				Ingredients: []ingredient.Ingredient{
					{Name: "Onion", MeasureType: "unit", Quantity: 1},
					{Name: "Rice", MeasureType: "mg", Quantity: 500},
					{Name: "Garlic", MeasureType: "unit", Quantity: 2},
				},
			},
			shouldExist: true,
		},
		{
			recipe: recipe.Recipe{
				Name: "Tomato Soup",
				Ingredients: []ingredient.Ingredient{
					{Name: "Tomato", MeasureType: "unit", Quantity: 4},
					{Name: "Water", MeasureType: "ml", Quantity: 500},
					{Name: "Salt", MeasureType: "mg", Quantity: 10},
				},
			},
			shouldExist: true,
		},
	}

	for _, tr := range testRecipes {
		if tr.shouldExist {
			err := service.AddRecipe(tr.recipe)
			if err != nil {
				t.Fatalf("failed to insert test recipe %s: %v", tr.recipe.Name, err)
			}
		}
	}

	t.Run("should retrieve all recipes with their ingredients", func(t *testing.T) {
		recipes, err := service.GetAllRecipes()
		if err != nil {
			t.Fatalf("failed to get all recipes: %v", err)
		}

		assert.Equal(t, len(testRecipes), len(recipes))

		for _, expectedRecipe := range testRecipes {
			if !expectedRecipe.shouldExist {
				continue
			}

			found := false
			for _, actualRecipe := range recipes {
				if actualRecipe.Name == expectedRecipe.recipe.Name {
					found = true
					assert.Equal(t, expectedRecipe.recipe.Name, actualRecipe.Name)
					assert.ElementsMatch(t, expectedRecipe.recipe.Ingredients, actualRecipe.Ingredients)
					break
				}
			}
			assert.True(t, found, "expected recipe %s not found", expectedRecipe.recipe.Name)
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
