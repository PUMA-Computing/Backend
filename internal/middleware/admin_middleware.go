package middleware

import (
	"Backend/internal/app/domain"
	"Backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := utils.GetUserRoleFromContext(c)

		if userRole != domain.RolePUMA {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Access Denied"})
		}

		return c.Next()
	}
}
