package v1

import (
	"Backend/internal/app/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupNewsRoutes(app *fiber.App, newsHandlers *handlers.UserHandlers) {
	//api := app.Group("/api/v1")
	//
	////api.Use(middleware.AuthMiddleware())
	////api.Post("/news/create", middleware.AdminMiddleware(), newsHandlers.CreateNews())
	////api.Put("/news/:id/edit", middleware.AdminMiddleware(), newsHandlers.EditNews())
	////api.Delete("/news/:id/delete", middleware.AdminMiddleware(), newsHandlers.DeleteNews())
	////
	////api.Get("/news", newsHandlers.GetNews())
	////api.Get("/news/:id", newsHandlers.GetNewsByID())
}
