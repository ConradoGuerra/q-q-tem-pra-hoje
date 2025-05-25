package main

import (
	"fmt"
	"net/http"
	"q-q-tem-pra-hoje/internal/app"
	"q-q-tem-pra-hoje/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("error loading .env file: %v", err)
	}

	connStr := database.SetupDB()
	db, err := database.Connect(connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	server := app.NewServer(db)

	fmt.Println("Server started at :8080")
	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
