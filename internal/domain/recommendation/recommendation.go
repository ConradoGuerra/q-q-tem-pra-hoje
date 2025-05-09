package recommendation

import "q-q-tem-pra-hoje/internal/domain/recipe"

type Recommendation struct {
	Recommendation int
	Recipe         recipe.Recipe
}
