package getUserContext

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const userIDKey = "userID"
const userRoleKey = "userRole"

func SetUserIDInContext(c *fiber.Ctx, userID uuid.UUID) {
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

func GetUserRoleFromContext(c *fiber.Ctx) (int, error) {
	userRole, ok := c.Locals(userRoleKey).(int)
	if !ok {
		return 0, errors.New("error getting userRole from context")
	}
	return userRole, nil
}

func GetUserIDAndRoleFromContext(c *fiber.Ctx) (uuid.UUID, int, error) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return uuid.UUID{}, 0, err
	}

	userRole, err := GetUserRoleFromContext(c)
	if err != nil {
		return uuid.UUID{}, 0, err
	}

	return userID, userRole, nil
}
