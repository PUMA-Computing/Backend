package getUserContext

import (
	"errors"
	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const userIDKey = "userID"
const userRoleKey = "userRole"

func SetUserIDInContext(c *fiber.Ctx, userID gocql.UUID) {
	c.Locals(userIDKey, userID)
}

func SetUserRoleInContext(c *fiber.Ctx, userRole string) {
	c.Locals(userRoleKey, userRole)
}

func GetUserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	userID, ok := c.Locals(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("error getting userID from context")
	}
	return userID, nil
}

func GetUserRoleFromContext(c *fiber.Ctx) (string, error) {
	userRole, ok := c.Locals(userRoleKey).(string)
	if !ok {
		return "", errors.New("error getting userRole from context")
	}
	return userRole, nil
}

func GetUserIDAndRoleFromContext(c *fiber.Ctx) (uuid.UUID, string, error) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return uuid.UUID{}, "", err
	}

	userRole, err := GetUserRoleFromContext(c)
	if err != nil {
		return uuid.UUID{}, "", err
	}

	return userID, userRole, nil
}
