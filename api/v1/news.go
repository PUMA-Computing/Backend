package v1

import (
	"Backend/internal/app/handlers/newsHandlers"
	"Backend/internal/middleware/admin"
	"github.com/gofiber/fiber/v2"
)

func NewsRoutes(app *fiber.App, newsHandlers *newsHandlers.NewsHandlers) {
	api := app.Group("/api/v1/news")

	api.Get("/", newsHandlers.GetNews())
	api.Get("/:id", newsHandlers.GetNewsByID())

	//api.Use(auth.AuthMiddleware())
	//api.Post("/:id/register", newsHandlers.LikeNews())

	api.Use(admin.AdminMiddleware())
	api.Post("/create", newsHandlers.CreateNews())
	api.Put("/:id/edit", newsHandlers.EditNews())
	api.Delete("/:id/delete", newsHandlers.DeleteNews())
}
