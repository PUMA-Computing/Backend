package v1

import (
	"Backend/internal/app/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, userHandlers *handlers.UserHandlers) {
	api := app.Group("/api/v1")

	api.Post("/users/register", userHandlers.RegisterUser())
	api.Post("/users/login", userHandlers.Login())
}
