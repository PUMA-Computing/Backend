package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserRoleFromContext(c *fiber.Ctx) string {
	token := c.Get("Authorization")

	claims := jwt.MapClaims{}
	_, _, err := new(jwt.Parser).ParseUnverified(token, claims)
	if err != nil {
		return ""
	}

	userRole, ok := claims["role"].(string)
	if !ok {
		return ""
	}

	return userRole
}
