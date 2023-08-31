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
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	cassandraRepo, err := repository2.NewCassandraRepository()
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

	err = app.Listen(":3000")
	if err != nil {
		return
	}
}
