package v1

import (
	"Backend/internal/app/handlers"
	"Backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupEventRoutes(app *fiber.App, eventHandlers *handlers.EventHandlers) {
	api := app.Group("/api/v1")

	api.Use(middleware.AuthMiddleware())

	api.Post("/event/create", middleware.AdminMiddleware(), eventHandlers.CreateEvent())
	api.Put("/event/:id/edit", middleware.AdminMiddleware(), eventHandlers.EditEvent())
	api.Delete("/event/:id/delete", middleware.AdminMiddleware(), eventHandlers.DeleteEvent())
	api.Get("/event/:id/users", middleware.AdminMiddleware(), eventHandlers.GetEventUsers())
	api.Post("/event/:id/register", middleware.AuthMiddleware(), eventHandlers.RegisterUserForEvent())

	api.Get("/event", eventHandlers.GetEvent())
	api.Get("/event/:id", eventHandlers.GetEventByID())
}
