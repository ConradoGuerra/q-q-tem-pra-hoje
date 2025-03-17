package postgres

import (
	"database/sql"
	"fmt"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
)

type ingredientStorageManager struct {
	db *sql.DB
}

func NewIngredientStorageManager(db *sql.DB) ingredientStorageManager {
	return ingredientStorageManager{db}
}

func (ism *ingredientStorageManager) AddIngredient(ingredient ingredient.Ingredient) error {
	query := "INSERT INTO ingredients_storage (name, measure_type, quantity) VALUES ($1, $2, $3)"
	_, err := ism.db.Exec(query, ingredient.Name, ingredient.MeasureType, ingredient.Quantity)
	if err != nil {
		return fmt.Errorf("failed to add ingredient: %v", err)
	}
	return nil
}

func (ism *ingredientStorageManager) FindIngredients() ([]ingredient.Ingredient, error) {
	query := "SELECT name, measure_type, sum(quantity) as quantity FROM ingredients_storage GROUP BY name, measure_type;"
	rows, err := ism.db.Query(query)

	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var ingredients []ingredient.Ingredient

	for rows.Next() {
		var ingredient ingredient.Ingredient

		err := rows.Scan(&ingredient.Name, &ingredient.MeasureType, &ingredient.Quantity)

		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		ingredients = append(ingredients, ingredient)
	}
	return ingredients, nil
}
