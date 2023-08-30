package middleware

import (
	"Backend/internal/app/domain"
	"Backend/internal/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionToken := c.Cookies("session_token")
		if sessionToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
		}

		userID, userRole, err := utils.ValidateSessionToken(sessionToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
		}

		if userRole != domain.RoleUser && userRole != domain.RolePUMA {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Access Denied"})
		}

		isValidSession, _ := utils.IsValidSessionToken(userID, sessionToken)
		if !isValidSession {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
		}

		if utils.IsTokenAboutToExpire(sessionToken, 5*time.Minute) {
			newToken, _ := utils.GenerateJWTToken(userID, domain.Role(userRole))
			c.Cookie(&fiber.Cookie{
				Name:     "session_token",
				Value:    newToken,
				Expires:  time.Now().Add(time.Hour * 24),
				Path:     "/",
				Secure:   true,
				HTTPOnly: true,
			})
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
