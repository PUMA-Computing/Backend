package main

import (
	v1 "Backend/api/v1"
	user3 "Backend/internal/app/handlers/user"
	repository2 "Backend/internal/app/interfaces/repository/cassandra"
	"Backend/internal/app/interfaces/repository/user"
	user2 "Backend/internal/app/interfaces/service/user"
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

	userRepository := user.NewCassandraUserRepository(cassandraRepo.Session)
	userService := user2.NewUserService(userRepository)
	userHandlers := user3.NewUserHandlers(userService)

	err = migrations.ExecuteMigrations(cassandraRepo.Session)
	if err != nil {
		log.Fatalf("Error executing migrations: %s", err.Error())
	}

	app := fiber.New()
	v1.AuthRoutes(app, userHandlers)
	err = app.Listen(":3000")
	if err != nil {
		return
	}
}
