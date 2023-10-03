package main

import (
	"Backend/api"
	"Backend/configs"
	"Backend/internal/database"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	config := configs.LoadConfig()

	database.Migrate()

	r := api.SetupRoutes()

	port := ":" + config.ServerPort
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
