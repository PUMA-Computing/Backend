package middleware

import (
	"Backend/internal/app/domain"
	"Backend/internal/utils"
	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
	"time"
)

func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionToken := c.Cookies("session_token")
		if sessionToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
		}

		userID, userRole, err := utils.ValidateSessionToken(sessionToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
		}

		userUUID, err := gocql.ParseUUID(userID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
		}

		if userRole != domain.RolePUMA {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Access Denied"})
		}

		isValidSession, _ := utils.IsValidSessionToken(userID, sessionToken)
		if !isValidSession {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
		}

		if utils.IsTokenAboutToExpire(sessionToken, 5*time.Minute) {
			newToken, _ := utils.GenerateJWTToken(userUUID, userRole)
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
