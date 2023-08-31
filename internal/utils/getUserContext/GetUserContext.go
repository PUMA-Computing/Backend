package getUserContext

import (
	"errors"
	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
)

const userIDKey = "userID"
const userRoleKey = "userRole"

func SetUserIDInContext(c *fiber.Ctx, userID gocql.UUID) {
	c.Locals(userIDKey, userID)
}

func SetUserRoleInContext(c *fiber.Ctx, userRole string) {
	c.Locals(userRoleKey, userRole)
}

func GetUserIDFromContext(c *fiber.Ctx) (gocql.UUID, error) {
	userID, ok := c.Locals(userIDKey).(gocql.UUID)
	if !ok {
		return gocql.UUID{}, errors.New("error getting userID from context")
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

func GetUserIDAndRoleFromContext(c *fiber.Ctx) (gocql.UUID, string, error) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return gocql.UUID{}, "", err
	}

	userRole, err := GetUserRoleFromContext(c)
	if err != nil {
		return gocql.UUID{}, "", err
	}

	return userID, userRole, nil
}
