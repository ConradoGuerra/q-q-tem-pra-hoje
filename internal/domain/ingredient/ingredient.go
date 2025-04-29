package ingredient

type Ingredient struct {
	ID          *int
	Name        string
	MeasureType string
	Quantity    int
}

func NewIngredient(id int, name string, measureType string, quantity int) Ingredient {
  ingredient := Ingredient{ID: &id, Name: name, MeasureType: measureType, Quantity: quantity}
	return ingredient
}
