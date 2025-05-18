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

func (ism *ingredientStorageManager) AddIngredient(ingredientParams ingredient.Ingredient) error {
	query := "SELECT id, quantity FROM ingredients_storage WHERE name = $1;"

	var ingredientFound ingredient.Ingredient
	err := ism.db.QueryRow(query, ingredientParams.Name).Scan(&ingredientFound.Id, &ingredientFound.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			query := "INSERT INTO ingredients_storage (name, measure_type, quantity) VALUES ($1, $2, $3)"
			_, err := ism.db.Exec(query, ingredientParams.Name, ingredientParams.MeasureType, ingredientParams.Quantity)
			if err != nil {
				return fmt.Errorf("failed to add ingredient: %v", err)
			}
			return nil
		}
		return fmt.Errorf("error executing query: %v", err)
	}

	newQuantity := ingredientFound.Quantity + ingredientParams.Quantity

	query = "UPDATE ingredients_storage SET quantity = $2 WHERE id = $1"

	_, err = ism.db.Exec(query, ingredientFound.Id, newQuantity)
	if err != nil {
		return fmt.Errorf("error to update ingredient: %v", err)
	}
	return nil

}

func (ism *ingredientStorageManager) FindIngredients() ([]ingredient.Ingredient, error) {
	query := "SELECT id, name, measure_type, quantity FROM ingredients_storage;"
	rows, err := ism.db.Query(query)

	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var ingredients []ingredient.Ingredient

	for rows.Next() {
		var ingredient ingredient.Ingredient

		err := rows.Scan(&ingredient.Id, &ingredient.Name, &ingredient.MeasureType, &ingredient.Quantity)

		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		ingredients = append(ingredients, ingredient)
	}
	return ingredients, nil
}

func (ism *ingredientStorageManager) Update(ingredientParams ingredient.Ingredient) error {
	query := "UPDATE ingredients_storage SET name = $1, quantity = $2, measure_type = $3 WHERE id = $4"
	_, err := ism.db.Exec(query, ingredientParams.Name, ingredientParams.Quantity, ingredientParams.MeasureType, ingredientParams.Id)
	if err != nil {
		return fmt.Errorf("error to update ingredient: %v", err)
	}
	return nil
}

func (ism *ingredientStorageManager) Delete(id uint) error {
	query := "DELETE from ingredients_storage WHERE id = $1"
	_, err := ism.db.Exec(query, id)

	if err != nil {
		return fmt.Errorf("error to delete ingredient: %v", err)
	}
	return nil
}
