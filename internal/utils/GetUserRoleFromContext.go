package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserRoleFromContext(c *fiber.Ctx) string {
	token := c.Get("Authorization")

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("pumacomputing"), nil
	})

	if err != nil {
		return ""
	}

	userRole := claims["role"].(string)
	return userRole
}
