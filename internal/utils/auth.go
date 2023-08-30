package utils

import "github.com/gofiber/fiber/v2"

func GetUserIDFromContext(c *fiber.Ctx) uint {
	userID := c.Locals("userID").(uint)
	return userID
}
