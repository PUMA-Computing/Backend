package main

import (
	"Backend/internal/app/handlers"
	"Backend/internal/app/repository"
	"Backend/internal/app/service"
	"github.com/gofiber/fiber/v2"
)

func main() {
	userRepository := repository.NewCassandraUserRepository(session)
	userService := service.NewUserService(userRepository)
	userHandlers := handlers.NewUserHandlers(userService)

	app := fiber.New()
	v1 := app.Group("/api/v1")
	v1users := v1.Group("/users")
	v1users.Post("/register", userHandlers.CreateUser())
	v1users.Post("/login", userHandlers.Login())
}
