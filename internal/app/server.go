package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"q-q-tem-pra-hoje/internal/database"
	"q-q-tem-pra-hoje/internal/repository/postgres"
	ingredientController "q-q-tem-pra-hoje/internal/server/controller/ingredient"
	recipeController "q-q-tem-pra-hoje/internal/server/controller/recipe"
	ingredientService "q-q-tem-pra-hoje/internal/service/ingredient"
	recipeService "q-q-tem-pra-hoje/internal/service/recipe"
)

type Server struct {
	server *http.Server
	db     *sql.DB
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewServer() (*Server, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	ism := postgres.NewIngredientStorageManager(db)
	rm := postgres.NewRecipeManager(db)
	is := ingredientService.NewService(&ism)
	rs := recipeService.NewRecipeService(rm)
	ic := ingredientController.NewIngredientController(is)
	rc := recipeController.NewRecipeController(is, rs)

	mux := http.NewServeMux()
	mux.Handle("/ingredient", ic)
	mux.Handle("/recipe", rc)

	handler := corsMiddleware(mux)

	return &Server{
		server: &http.Server{
			Addr:    ":8080",
			Handler: handler,
		},
		db: db,
	}, nil
}

func (s Server) Start() error {
	return s.server.ListenAndServe()
}

func (s Server) Close() error {
	return s.db.Close()
}
