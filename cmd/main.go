package main

import (
	v1 "Backend/api/v1"
	"Backend/internal/app/handlers"
	"Backend/internal/app/repository"
	"Backend/internal/app/service"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cassandraRepo, err := repository.NewCassandraRepository()
	if err != nil {
		panic(err)
	}
	defer cassandraRepo.Close()

	userRepository := repository.NewCassandraUserRepository(cassandraRepo.Session)
	userService := service.NewUserService(userRepository)
	userHandlers := handlers.NewUserHandlers(userService)

	app := fiber.New()
	v1.AuthRoutes(app, userHandlers)
	err = app.Listen(":3000")
	if err != nil {
		return
	}
}
