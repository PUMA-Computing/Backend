package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateSessionToken(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", "", errors.New("invalid token claims")
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		return "", "", errors.New("missing userID claim")
	}

	userRole, ok := claims["role"].(string)
	if !ok {
		return "", "", errors.New("missing role claim")
	}

	return userID, userRole, nil
}
