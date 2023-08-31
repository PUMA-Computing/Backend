package v1

import (
	"Backend/internal/app/handlers/event"
	"Backend/internal/middleware/admin"
	"Backend/internal/middleware/auth"
	"github.com/gofiber/fiber/v2"
)

func EventRoutes(app *fiber.App, eventHandlers *event.EventHandlers) {
	api := app.Group("/api/v1")

	api.Use(auth.AuthMiddleware())

	api.Post("/event/create", admin.AdminMiddleware(), eventHandlers.CreateEvent())
	api.Put("/event/:id/edit", admin.AdminMiddleware(), eventHandlers.EditEvent())
	api.Delete("/event/:id/delete", admin.AdminMiddleware(), eventHandlers.DeleteEvent())
	api.Get("/event/:id/users", admin.AdminMiddleware(), eventHandlers.GetEventUsers())
	api.Post("/event/:id/register", auth.AuthMiddleware(), eventHandlers.RegisterUserForEvent())

	api.Get("/event", eventHandlers.GetEvent())
	api.Get("/event/:id", eventHandlers.GetEventByID())
}
