package ingredient

type Ingredient struct {
	Name        string
	MeasureType string
	Quantity    int
}

func NewIngredient(name string, measureType string, quantity int) Ingredient {
	ingredient := Ingredient{Name: name, MeasureType: measureType, Quantity: quantity}
	return ingredient
}
