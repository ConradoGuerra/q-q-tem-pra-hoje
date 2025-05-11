package app

import (
	"database/sql"
	"net/http"
	"q-q-tem-pra-hoje/internal/repository/postgres"
	ingredientController "q-q-tem-pra-hoje/internal/server/controller/ingredient"
	recipeController "q-q-tem-pra-hoje/internal/server/controller/recipe"
	recommendationController "q-q-tem-pra-hoje/internal/server/controller/recommendation"
	ingredientService "q-q-tem-pra-hoje/internal/service/ingredient"
	recipeService "q-q-tem-pra-hoje/internal/service/recipe"
	recommendationService "q-q-tem-pra-hoje/internal/service/recommendation"
)

type Server struct {
	server *http.Server
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewServer(db *sql.DB) (*Server, error) {

	ism := postgres.NewIngredientStorageManager(db)
	rm := postgres.NewRecipeManager(db)
	is := ingredientService.NewService(&ism)
	rs := recipeService.NewRecipeService(rm)
	res := recommendationService.NewRecommendationService(rm)
	ic := ingredientController.NewIngredientController(is)
	rc := recipeController.NewRecipeController(is, rs)
	rec := recommendationController.NewRecommendationController(is, res)

	mux := http.NewServeMux()
	mux.Handle("/ingredient", ic)
	mux.Handle("/ingredient/{id}", ic)
	mux.Handle("/recipe", rc)
	mux.Handle("/recommendation", rec)

	handler := corsMiddleware(mux)

	return &Server{
		server: &http.Server{
			Addr:    ":8080",
			Handler: handler,
		}}, nil
}

func (s Server) Start() error {
	return s.server.ListenAndServe()
}
