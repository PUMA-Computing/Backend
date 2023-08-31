package v1

import (
	"Backend/internal/app/handlers/userHandlers"
	"github.com/gofiber/fiber/v2"
)

func SetupNewsRoutes(app *fiber.App, newsHandlers *userHandlers.UserHandlers) {
	//api := app.Group("/api/v1")
	//
	////api.Use(middleware.AuthMiddleware())
	////api.Post("/newsHandlers/create", middleware.AdminMiddleware(), newsHandlers.CreateNews())
	////api.Put("/newsHandlers/:id/edit", middleware.AdminMiddleware(), newsHandlers.EditNews())
	////api.Delete("/newsHandlers/:id/delete", middleware.AdminMiddleware(), newsHandlers.DeleteNews())
	////
	////api.Get("/newsHandlers", newsHandlers.GetNews())
	////api.Get("/newsHandlers/:id", newsHandlers.GetNewsByID())
}
