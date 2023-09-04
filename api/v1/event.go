package v1

import (
	"Backend/internal/app/handlers/eventHandlers"
	"Backend/internal/app/interfaces/repository/userRepository"
	"Backend/internal/middleware/admin"
	"Backend/internal/middleware/auth"
	"github.com/gofiber/fiber/v2"
)

func EventRoutes(app *fiber.App, eventHandlers *eventHandlers.EventHandlers, userRepo userRepository.UserRepository) {
	api := app.Group("/api/v1/event")

	api.Get("/", eventHandlers.GetEvent())
	api.Get("/:id", eventHandlers.GetEventByID())

	api.Use(auth.Middleware(userRepo))
	api.Post("/:id/register", eventHandlers.RegisterUserForEvent())

	api.Use(admin.Middleware(userRepo))
	api.Post("/create", eventHandlers.CreateEvent())
	api.Put("/:id/edit", eventHandlers.EditEvent())
	api.Delete("/:id/delete", eventHandlers.DeleteEvent())
	api.Get("/:id/users", eventHandlers.GetEventUsers())

}
