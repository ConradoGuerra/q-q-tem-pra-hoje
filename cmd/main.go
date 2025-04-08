package main

import (
	"fmt"
	"net/http"
	"q-q-tem-pra-hoje/internal/app"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("error loading .env file: %v", err)
	}

	server, err := app.NewServer()
	if err != nil {
		panic(err)
	}
	defer server.Close()
	fmt.Println("Server started at :8080")
	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
