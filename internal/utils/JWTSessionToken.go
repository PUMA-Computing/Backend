package utils

import (
	"Backend/internal/app/domain"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	jwtSecretKey          = "pumacomputing"
	jwtExpirationDuration = time.Hour * 24
)

func GenerateJWTToken(userID string, role domain.Role) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"role":   role,
		"exp":    time.Now().Add(jwtExpirationDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
