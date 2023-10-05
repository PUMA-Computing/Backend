package main

import (
	"Backend/api"
	"Backend/configs"
	"Backend/internal/database"
	"Backend/pkg/utils"
	"github.com/gin-contrib/cors"
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

	r := api.SetupRoutes()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}

	r.Use(cors.New(corsConfig))

	port := ":" + config.ServerPort
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
