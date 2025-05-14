package recipe

import (
	"errors"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
)

type Recipe struct {
	Id          *int
	Name        string
	Ingredients []ingredient.Ingredient
}

func (r *Recipe) Validate() error {
	if r.Name == "" {
		return errors.New("recipe name cannot be empty")
	}
	if len(r.Ingredients) == 0 {
		return errors.New("recipe must have at least one ingredient")
	}
	return nil
}

func NewRecipe(id int, name string, ingredients []ingredient.Ingredient) (Recipe, error) {
	recipe := Recipe{Id: &id, Name: name, Ingredients: ingredients}
	if err := recipe.Validate(); err != nil {
		return Recipe{}, err
	}
	return recipe, nil
}
