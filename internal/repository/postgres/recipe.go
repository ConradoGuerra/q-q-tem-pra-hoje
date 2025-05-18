package postgres

import (
	"database/sql"
	"fmt"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recipe"
)

type recipeManager struct {
	*sql.DB
}

func NewRecipeManager(db *sql.DB) *recipeManager {
	return &recipeManager{db}
}

func (rm recipeManager) AddRecipe(recipe recipe.Recipe) error {
	var recipeId  int
	err := rm.QueryRow("INSERT INTO recipes (name) VALUES ($1) RETURNING id;", recipe.Name).Scan(&recipeId )
	if err != nil {
		return fmt.Errorf("failed to insert recipe: %v", err)
	}

	for _, ing := range recipe.Ingredients {
		_, err = rm.Exec(`
		      INSERT INTO recipes_ingredients (recipe_id, name, measure_type, quantity)
		      VALUES ($1, $2, $3, $4)
		      ON CONFLICT (recipe_id, name) DO NOTHING;
		  `, recipeId , ing.Name, ing.MeasureType, ing.Quantity)
		if err != nil {
			return fmt.Errorf("failed to insert a recipe ingredient: %v", err)
		}
	}
	return nil
}
func (rm recipeManager) GetAllRecipes() ([]recipe.Recipe, error) {
	rows, err := rm.Query(`SELECT 
                          r.id,
                          r.name, 
                          i.name, 
                          i.measure_type, 
                          i.quantity 
                        FROM recipes r 
                          JOIN recipes_ingredients i ON r.id = i.recipe_id`)

	if err != nil {
		return nil, fmt.Errorf("error querying recipes: %w", err)
	}

	defer rows.Close()

	recipeMap := make(map[string]*recipe.Recipe)

	for rows.Next() {
		var recipeId int
		var recipeName string
		var ingredientRetrieved ingredient.Ingredient

		err := rows.Scan(&recipeId, &recipeName, &ingredientRetrieved.Name, &ingredientRetrieved.MeasureType, &ingredientRetrieved.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		if r, exists := recipeMap[recipeName]; exists {
			r.Ingredients = append(r.Ingredients, ingredientRetrieved)
		} else {
			newRecipe := &recipe.Recipe{
				Id:          &recipeId,
				Name:        recipeName,
				Ingredients: []ingredient.Ingredient{ingredientRetrieved},
			}
			recipeMap[recipeName] = newRecipe
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	var recipesRetrieved []recipe.Recipe
	for _, r := range recipeMap {
		recipesRetrieved = append(recipesRetrieved, *r)
	}

	return recipesRetrieved, nil
}

func (rm recipeManager) DeleteRecipe(id uint) error {
	_, err := rm.Exec(`
		DELETE FROM recipes_ingredients 
		WHERE recipe_id IN (SELECT id FROM recipes WHERE id = $1)
	`, id)
	if err != nil {
		return fmt.Errorf("failed to delete recipe ingredients: %v", err)
	}

	result, err := rm.Exec("DELETE FROM recipes WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete recipe: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("recipe not found")
	}

	return nil
}
