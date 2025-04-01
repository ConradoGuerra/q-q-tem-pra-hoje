package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"q-q-tem-pra-hoje/internal/database"
	"q-q-tem-pra-hoje/internal/repository/postgres"
	controller "q-q-tem-pra-hoje/internal/server/controller/ingredient"
	service "q-q-tem-pra-hoje/internal/service/ingredient"
)

type Server struct {
	server *http.Server
	db     *sql.DB
}

func NewServer() (*Server, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}
	defer db.Close()

	manager := postgres.NewIngredientStorageManager(db)
	service := service.NewService(&manager)
	ingredientController := controller.NewIngredientController(service)

	mux := http.NewServeMux()
	mux.Handle("/ingredient", ingredientController)

	return &Server{
		server: &http.Server{
			Addr:    ":8080",
			Handler: mux,
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
