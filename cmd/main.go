//   PUMA Computing API:
//    version: 1.0.0
//    title: PUMA Computing API
//   Schemes: http, https
//   Host: localhost:3000
//   BasePath: /api/v1
//      Consumes:
//      - application/json
//   Produces:
//   - application/json
//   SecurityDefinitions:
//    Bearer:
//     type: apiKey
//     name: Authorization
//     in: header
//   swagger:meta
package main

import (
	v1 "Backend/api/v1"
	eventHandlers2 "Backend/internal/app/handlers/eventHandlers"
	news3 "Backend/internal/app/handlers/newsHandlers"
	user3 "Backend/internal/app/handlers/userHandlers"
	event1 "Backend/internal/app/interfaces/repository/eventRepository"
	news1 "Backend/internal/app/interfaces/repository/newsRepository"
	"Backend/internal/app/interfaces/repository/postgresRepository"
	user1 "Backend/internal/app/interfaces/repository/userRepository"
	event2 "Backend/internal/app/interfaces/service/eventService"
	news2 "Backend/internal/app/interfaces/service/newsService"
	user2 "Backend/internal/app/interfaces/service/userService"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

/**
 * Connect to database with retry
 * This function will try to connect to the database 5 times with a 5s delay between each try.
 */
func connectToDatabaseWithRetry() (*postgresRepository.PostgresRepository, error) {
	var postgresRepo *postgresRepository.PostgresRepository
	var err error
	for i := 0; i < 5; i++ {
		postgresRepo, err = postgresRepository.NewPostgresRepository()
		if err == nil {
			return postgresRepo, nil
		}
		time.Sleep(5 * time.Second)
	}
	return nil, err
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	jwtSecretKey := os.Getenv("JWT_SECRET")
	if jwtSecretKey == "" {
		log.Fatalf("JWT_SECRET_KEY environment variable not set")
	}

	postgresRepo, err := connectToDatabaseWithRetry()
	if err != nil {
		panic(err)
	}
	defer postgresRepo.Close()

	/**
	 * Repositories
	 */
	userRepository := user1.NewPostgresUserRepository(postgresRepo.DB)
	eventRepository := event1.NewPostgresForEventRepository(postgresRepo.DB)
	newsRepository := news1.NewPostgresForNewsRepository(postgresRepo.DB)

	/**
	 * Services
	 */
	userService := user2.NewUserService(userRepository)
	eventService := event2.NewEventService(eventRepository)
	newsService := news2.NewNewsService(newsRepository)

	/**
	 * Handlers
	 */
	userHandlers := user3.NewUserHandlers(userService)
	eventHandlers := eventHandlers2.NewEventHandlers(eventService)
	newsHandlers := news3.NewNewsHandlers(newsService)

	app := fiber.New()

	/**
	 * Routes
	 */
	v1.AuthRoutes(app, userHandlers)
	v1.EventRoutes(app, eventHandlers)
	v1.NewsRoutes(app, newsHandlers)

	err = app.Listen(":3000")
	if err != nil {
		return
	}
}
