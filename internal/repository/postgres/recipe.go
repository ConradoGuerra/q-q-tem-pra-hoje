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
	var recipeID int
	err := rm.QueryRow("INSERT INTO recipes (name) VALUES ($1) RETURNING id;", recipe.Name).Scan(&recipeID)
	if err != nil {
		return fmt.Errorf("failed to insert recipe: %v", err)
	}

	for _, ing := range recipe.Ingredients {
		_, err = rm.Exec(`
		      INSERT INTO recipes_ingredients (recipe_id, name, measure_type, quantity)
		      VALUES ($1, $2, $3, $4)
		      ON CONFLICT (recipe_id, name) DO NOTHING;
		  `, recipeID, ing.Name, ing.MeasureType, ing.Quantity)
		if err != nil {
			return fmt.Errorf("failed to insert a recipe ingredient: %v", err)
		}
	}
	return nil
}
func (rm recipeManager) GetAllRecipes() ([]recipe.Recipe, error) {
	rows, err := rm.Query(`SELECT 
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
		var recipeName string
		var ingredientRetrieved ingredient.Ingredient

		err := rows.Scan(&recipeName, &ingredientRetrieved.Name, &ingredientRetrieved.MeasureType, &ingredientRetrieved.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		if r, exists := recipeMap[recipeName]; exists {
			r.Ingredients = append(r.Ingredients, ingredientRetrieved)
		} else {
			newRecipe := &recipe.Recipe{
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
