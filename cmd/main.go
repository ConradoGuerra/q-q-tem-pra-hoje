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

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("error loading .env file: %v", err)
	}

	db, err := database.Connect()
	if err != nil {
		fmt.Printf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	manager := postgres.NewIngredientStorageManager(db)

	service := service.NewService(&manager)

	ingredientController := controller.NewIngredientController(service)

	http.HandleFunc("/ingredient", ingredientController.Add)
	fmt.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Could not start server: %v\n", err)
	}
}
