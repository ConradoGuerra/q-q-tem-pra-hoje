package ingredient

type Ingredient struct {
	Id          *int
	Name        string
	MeasureType string
	Quantity    int
}

func NewIngredient(id int, name string, measureType string, quantity int) Ingredient {
  ingredient := Ingredient{Id: &id, Name: name, MeasureType: measureType, Quantity: quantity}
	return ingredient
}
