package v1

import (
	"Backend/internal/app/handlers/userHandlers"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, userHandlers *userHandlers.UserHandlers) {
	api := app.Group("/api/v1")

	api.Post("/auth/register", userHandlers.RegisterUser())
	api.Post("/auth/login", userHandlers.Login())
	api.Post("/auth/logout", userHandlers.Logout())
}
