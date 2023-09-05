package auth

import (
	"Backend/internal/app/domain/roles"
	"Backend/internal/app/interfaces/repository/userRepository"
	token2 "Backend/internal/utils/token"
	"fmt"
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

		authToken := token2.ExtractBearerToken(authHeader)
		if authToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Extract Token Failed"})
		}

		userID, err := token2.ValidateSessionToken(authToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Session Token Not Valid", "error": err.Error()})
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized", "error": err.Error()})
		}

		userRoleID, err := userRepo.GetUserRoleByID(userUUID)
		//fmt.Printf("userRoleID: %d\n", userRoleID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized", "error": err.Error()})
		}

		if userRoleID != roles.RolePUMA && userRoleID != roles.RoleComputizen {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized", "error": err.Error()})
		}

		isValid, err := token2.IsValidSessionToken(userID, authToken)
		if err != nil || !isValid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid session token", "error": err.Error()})
		}

		if token2.IsSessionTokenAboutExpired(authToken, 5*time.Minute) {
			newToken, err := token2.GenerateJWTToken(userUUID, userRoleID)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid session token", "error": err.Error()})
			}

			c.Response().Header.Set("Authorization", "Bearer "+newToken)
		}
		c.Locals("userID", userUUID)

		return c.Next()
	}
}

func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println("Request: %s &s\n", c.Method(), c.Path())
		return c.Next()
	}
}
