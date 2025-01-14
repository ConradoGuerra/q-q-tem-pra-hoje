package ingredient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Ingredient struct {
	Name        string
	MeasureType string
	Quantity    int
}

// Interface
type IngredientManager interface {
	AddIngredient(Ingredient)
	FindIngredients() []Ingredient
}

// Service
type IngredientInventoryService struct {
	ingredientRepository IngredientManager
}

// Instance
func NewService(ingredientRepository IngredientManager) *IngredientInventoryService {
	return &IngredientInventoryService{ingredientRepository}
}

// Implements method
func (i *IngredientInventoryService) AddIngredientToInventory(ingredient Ingredient) {
	i.ingredientRepository.AddIngredient(ingredient)
}

func (i *IngredientInventoryService) FindIngredients() []Ingredient {
	return i.ingredientRepository.FindIngredients()
}

// Repo
type InMemoryIngredientRepository struct {
	ingredients []Ingredient
}

// Implement interface
func (r *InMemoryIngredientRepository) AddIngredient(ingredient Ingredient) {
	r.ingredients = append(r.ingredients, ingredient)
}
func (r *InMemoryIngredientRepository) FindIngredients() []Ingredient {
	ingredientMap := make(map[string]Ingredient)
	for _, ingredient := range r.ingredients {
		if existing, exists := ingredientMap[ingredient.Name]; exists {
			existing.Quantity += ingredient.Quantity
			ingredientMap[ingredient.Name] = existing
		} else {
			ingredientMap[ingredient.Name] = ingredient
		}
	}

	ingredientsFound := make([]Ingredient, 0, len(ingredientMap))
	for _, ingredient := range ingredientMap {
		ingredientsFound = append(ingredientsFound, ingredient)
	}
	return ingredientsFound
}

func TestAddIngredientToInventory(t *testing.T) {
	t.Run("it should add ingredients to inventory", func(t *testing.T) {

		repository := &InMemoryIngredientRepository{}
		ingredientService := NewService(repository)

		ingredient := Ingredient{Name: "onion", Quantity: 10, MeasureType: "unit"}
		ingredientService.AddIngredientToInventory(ingredient)

		assert.Contains(t, repository.ingredients, ingredient, "Ingredient should be added to inventory")
	})
}
func TestFindIngredients(t *testing.T) {

	t.Run("it should find all ingredients in the inventory", func(t *testing.T) {
		repository := &InMemoryIngredientRepository{}
		ingredientService := NewService(repository)
		ingredientService.AddIngredientToInventory(Ingredient{Name: "onion", Quantity: 10, MeasureType: "unit"})
		ingredientService.AddIngredientToInventory(Ingredient{Name: "garlic", Quantity: 2, MeasureType: "unit"})
		ingredientService.AddIngredientToInventory(Ingredient{Name: "onion", Quantity: 10, MeasureType: "unit"})
		ingredients := ingredientService.FindIngredients()
		expectedIngredients := []Ingredient{{Name: "onion", Quantity: 20, MeasureType: "unit"}, {Name: "garlic", Quantity: 2, MeasureType: "unit"}}
		assert.Equal(t, expectedIngredients, ingredients)
	})
}
