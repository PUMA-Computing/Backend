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
	repository2 "Backend/internal/app/interfaces/repository/cassandraRepository"
	event1 "Backend/internal/app/interfaces/repository/eventRepository"
	news1 "Backend/internal/app/interfaces/repository/newsRepository"
	user1 "Backend/internal/app/interfaces/repository/userRepository"
	event2 "Backend/internal/app/interfaces/service/eventService"
	news2 "Backend/internal/app/interfaces/service/newsService"
	user2 "Backend/internal/app/interfaces/service/userService"
	"Backend/pkg/migrations"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

/**
 * Connect to database with retry
 * This function will try to connect to the database 5 times with a 5s delay between each try.
 */
func connectToDatabaseWithRetry() (*repository2.CassandraRepository, error) {
	var cassandraRepo *repository2.CassandraRepository
	var err error

	for retries := 0; retries < 5; retries++ {
		cassandraRepo, err = repository2.NewCassandraRepository()
		if err == nil {
			break
		}
		log.Printf("Error connecting to the database: %s. Retrying...", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return nil, err
	}
	return cassandraRepo, nil
}

func main() {
	cassandraRepo, err := connectToDatabaseWithRetry()
	if err != nil {
		panic(err)
	}
	defer cassandraRepo.Close()

	/**
	 * Repositories
	 */
	userRepository := user1.NewCassandraUserRepository(cassandraRepo.Session)
	eventRepository := event1.NewCassandraForEventRepository(cassandraRepo.Session)
	newsRepository := news1.NewCassandraForNewsRepository(cassandraRepo.Session)

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

	err = migrations.ExecuteMigrations(cassandraRepo.Session)
	if err != nil {
		log.Fatalf("Error executing migrations: %s", err.Error())
	}

	app := fiber.New()

	/**
	 * Routes
	 */
	v1.AuthRoutes(app, userHandlers)
	v1.EventRoutes(app, eventHandlers)
	v1.NewsRoutes(app, newsHandlers)

	/**
	 * Database connection checker
	 *
	 * This goroutine will check if the database connection is still alive every 10 seconds.
	 */
	go func() {
		for {
			if cassandraRepo.Session.Closed() {
				log.Println("Database connection lost. Reconnecting")
				cassandraRepo, err = connectToDatabaseWithRetry()
				if err != nil {
					log.Printf("Failed to reconnect to database: %s", err.Error())
				} else {
					log.Println("Reconnected to database")
				}
			}
			time.Sleep(10 * time.Second)
		}
	}()

	err = app.Listen(":3000")
	if err != nil {
		return
	}
}
