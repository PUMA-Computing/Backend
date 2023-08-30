package middleware

import (
	"Backend/internal/app/domain"
	"Backend/internal/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := utils.GetUserRoleFromContext(c)

		if userRole != domain.RoleUser && userRole != domain.RolePUMA {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Access Denied"})
		}
		return c.Next()
	}
}

func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println("Request: %s &s\n", c.Method(), c.Path())
		return c.Next()
	}
}
