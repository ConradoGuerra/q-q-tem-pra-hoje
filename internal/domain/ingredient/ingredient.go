package ingredient

import (
	"errors"
)

type Ingredient struct {
	Name        string
	MeasureType string
	Quantity    int
}

func (i *Ingredient) Validate() []error {
	var invalidations []error
	if i.Name == "" {
		invalidations = append(invalidations, errors.New("ingredient must be a not empty string"))
	}
	if i.MeasureType == "" {
		invalidations = append(invalidations, errors.New("measure_type must be a not empty string"))
	}
	if len(invalidations) != 0 {
		return invalidations
	}
	return nil
}

func NewIngredient(name string, measureType string, quantity int) (Ingredient, []error) {
	ingredient := Ingredient{Name: name, MeasureType: measureType, Quantity: quantity}
	if err := ingredient.Validate(); err != nil {
		return Ingredient{}, err
	}
	return ingredient, nil
}
