package auth

import (
	"Backend/internal/app/domain/roles"
	"Backend/internal/utils/token"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionToken := c.Cookies("session_token")
		if sessionToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
		}

		userID, userRole, err := token.ValidateSessionToken(sessionToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
		}

		if userRole != roles.RoleComputizen && userRole != roles.RolePUMA {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Access Denied"})
		}

		isValidSession, _ := token.IsValidSessionToken(userID, sessionToken)
		if !isValidSession {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
		}

		if token.IsTokenAboutToExpire(sessionToken, 5*time.Minute) {
			newToken, _ := token.GenerateJWTToken(userUUID, userRole)
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
