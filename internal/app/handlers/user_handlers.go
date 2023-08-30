package handlers

import (
	"Backend/internal/app/domain"
	"Backend/internal/app/service"
	"github.com/gofiber/fiber/v2"
)

type UserHandlers struct {
	userService service.UserServices
}

func NewUserHandlers(userService service.UserServices) *UserHandlers {
	return &UserHandlers{userService: userService}
}

func (h *UserHandlers) CreateUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user domain.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error parsing user"})
		}

		if err := h.userService.CreateUser(&user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error creating user"})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created successfully"})
	}
}

func (h *UserHandlers) Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var loginData struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.BodyParser(&loginData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error parsing login data"})
		}

		user, err := h.userService.AuthenticateUser(loginData.Email, loginData.Password)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid email or password"})
		}

		return c.Status(fiber.StatusOK).JSON(user)
	}
}
