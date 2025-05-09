package recommendation

import "q-q-tem-pra-hoje/internal/domain/ingredient"

type RecommendationProvider interface{
	GetRecommendations(ingredients *[]ingredient.Ingredient) ([]Recommendation, error)

}
