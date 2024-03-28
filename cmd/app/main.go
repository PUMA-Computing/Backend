package main

import (
	"Backend/api"
	"Backend/configs"
	"Backend/internal/database"
	"Backend/internal/services"
	"Backend/pkg/utils"
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
	database.Init(config)
	utils.InitRedis()

	eventService := services.NewEventService()
	eventStatusUpdater := services.NewEventStatusUpdater(eventService)
	go eventStatusUpdater.Run()

	r := api.SetupRoutes()

	port := ":"
	if config.ApiPort == "" {
		log.Fatalf("API port is not set")
	} else {
		port += config.ApiPort
	}
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
