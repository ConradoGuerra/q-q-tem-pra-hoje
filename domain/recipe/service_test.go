package recipe

import (
	"errors"
	"q-q-tem-pra-hoje/domain/ingredient"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Recipe struct {
	Name        string
	Ingredients []ingredient.Ingredient
}

type RecipeManager interface {
	AddRecipe(recipe Recipe) error
}

func NewRecipe(name string, ingredients []ingredient.Ingredient) (Recipe, error) {
	recipe := Recipe{Name: name, Ingredients: ingredients}
	if err := recipe.Validate(); err != nil {
		return Recipe{}, err
	}
	return recipe, nil
}

type RecipeService struct {
	RecipeManager
}

func NewRecipeService(m RecipeManager) *RecipeService {
	return &RecipeService{RecipeManager: m}
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

func (s *RecipeService) AddRecipe(recipe Recipe) error {
	return s.RecipeManager.AddRecipe(recipe)
}

type InMemoryRecipeManager struct {
	Recipes []Recipe
}

func NewInMemoryRecipeManager() *InMemoryRecipeManager {
	return &InMemoryRecipeManager{Recipes: []Recipe{}}
}

func (m *InMemoryRecipeManager) AddRecipe(recipe Recipe) error {
	m.Recipes = append(m.Recipes, recipe)
	return nil
}

func TestAddRecipe(t *testing.T) {
	t.Run("It should add a valid recipe", func(t *testing.T) {

		expectedRecipe, err := NewRecipe("Rice", []ingredient.Ingredient{
			{Name: "Onion", MeasureType: "unit", Quantity: 1},
			{Name: "Rice", MeasureType: "mg", Quantity: 500},
			{Name: "Garlic", MeasureType: "unit", Quantity: 2}})
		inMemoryRecipeManager := NewInMemoryRecipeManager()
		recipeService := NewRecipeService(inMemoryRecipeManager)

		recipeService.AddRecipe(expectedRecipe)
		assert.NoError(t, err)
		assert.Equal(t, expectedRecipe, inMemoryRecipeManager.Recipes[len(inMemoryRecipeManager.Recipes)-1])
	})

	t.Run("It should return an error for an invalid name", func(t *testing.T) {
		invalidRecipe, err := NewRecipe("", []ingredient.Ingredient{ // Empty name
			{Name: "Onion", MeasureType: "unit", Quantity: 1},
			{Name: "Rice", MeasureType: "mg", Quantity: 500},
			{Name: "Garlic", MeasureType: "unit", Quantity: 2},
		})
		manager := NewInMemoryRecipeManager()
		service := NewRecipeService(manager)

		service.AddRecipe(invalidRecipe)
		assert.Error(t, err)
		assert.Equal(t, "recipe name cannot be empty", err.Error())
		assert.Equal(t, Recipe{}, invalidRecipe)
	})

	t.Run("It should return an error for invalid ingredients", func(t *testing.T) {
		invalidRecipe, err := NewRecipe("Rice", []ingredient.Ingredient{}) // No ingredients
		manager := NewInMemoryRecipeManager()
		service := NewRecipeService(manager)

		service.AddRecipe(invalidRecipe)
		assert.Error(t, err)
		assert.Equal(t, "recipe must have at least one ingredient", err.Error())
		assert.Equal(t, Recipe{}, invalidRecipe)

	})
}
