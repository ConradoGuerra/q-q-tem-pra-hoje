package main

import (
	"fmt"
	"net/http"
	"q-q-tem-pra-hoje/internal/database"
	"q-q-tem-pra-hoje/internal/repository/postgres"
	controller "q-q-tem-pra-hoje/internal/server/controller/ingredient"
	service "q-q-tem-pra-hoje/internal/service/ingredient"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func setupServer() (*http.Server, error) {
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

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}, nil
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("error loading .env file: %v", err)
	}

	server, err := setupServer()
	if err != nil {
		server.ErrorLog.Panicf("Failed to setup server: %v", err)
		return
	}
	defer server.Close()
	fmt.Println("Server started at " + server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		server.ErrorLog.Panicf("Could not start server: %v\n", err)
	}
}
