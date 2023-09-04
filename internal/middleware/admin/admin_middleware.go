package admin

import (
	"Backend/internal/app/domain/roles"
	"Backend/internal/app/interfaces/repository/userRepository"
	"Backend/internal/utils/token"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

func Middleware(userRepo userRepository.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Missing Authorization header"})
		}

		authToken := token.ExtractBearerToken(authHeader)
		if authToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Extract Token Failed"})
		}

		userID, err := token.ValidateSessionToken(authToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Session Token Not Valid", "error": err.Error(), "Extract Token:": authToken})
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized", "error": err.Error()})
		}

		userRoleID, err := userRepo.GetUserRoleByID(userUUID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized", "error": err.Error()})
		}

		if userRoleID != roles.RolePUMA {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
		}

		isValidSession, _ := token.IsValidSessionToken(userID, authToken)
		if !isValidSession {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid session token", "error": err.Error()})
		}

		if token.IsSessionTokenAboutExpired(authToken, 5*time.Minute) {
			newToken, _ := token.GenerateJWTToken(userUUID, userRoleID)

			c.Response().Header.Set("Authorization", "Bearer "+newToken)
		}

		return c.Next()
	}
}
