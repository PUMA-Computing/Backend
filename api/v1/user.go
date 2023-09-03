package v1

import (
	"Backend/internal/app/handlers/userHandlers"
	"Backend/internal/middleware/admin"
	"Backend/internal/middleware/auth"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, userHandlers *userHandlers.UserHandlers) {
	api := app.Group("/api/v1")

	api.Post("/auth/register", userHandlers.RegisterUser())
	api.Post("/auth/login", userHandlers.Login())
	api.Post("/auth/logout", userHandlers.Logout())

	api.Get("/user/", admin.Middleware(), userHandlers.GetAllUsers())
	api.Get("/user/:id", auth.Middleware(), userHandlers.GetUserProfile())
	api.Post("/user/update", auth.Middleware(), userHandlers.UpdateUser())
	api.Post("/user/delete", auth.Middleware(), userHandlers.DeleteUser())
	//api.Post("/user/forgot-password", userHandlers.ForgotPassword())
}
