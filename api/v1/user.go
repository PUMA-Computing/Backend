package v1

import (
	"Backend/internal/app/handlers/userHandlers"
	"Backend/internal/app/interfaces/repository/userRepository"
	"Backend/internal/middleware/admin"
	"Backend/internal/middleware/auth"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, userHandlers *userHandlers.UserHandlers, userRepo userRepository.UserRepository) {
	api := app.Group("/api/v1")

	api.Post("/auth/register", userHandlers.RegisterUser())
	api.Post("/auth/login", userHandlers.Login())
	//api.Post("/user/forgot-password", userHandlers.ForgotPassword())

	api.Get("/user/:id", auth.Middleware(userRepo, "Manage_Profile"), userHandlers.GetUserProfile())
	api.Post("/user/update", userHandlers.UpdateUser())
	api.Post("/user/delete", userHandlers.DeleteUser())
	api.Post("/user/logout", userHandlers.Logout())

	api.Use(admin.Middleware(userRepo))
	api.Get("/user/", userHandlers.GetAllUsers())
}
