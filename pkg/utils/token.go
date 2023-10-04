package utils

import (
	"context"
	"github.com/form3tech-oss/jwt-go"
	"github.com/google/uuid"
	"time"
)

func GenerateJWTToken(userID uuid.UUID, secretKey string) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type CustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}

func StoreTokenInRedis(userID uuid.UUID, token string) error {
	ctx := context.Background()
	ttl := 24 * time.Hour
	err := Rdb.Set(ctx, userID.String(), token, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func RetrieveTokenFromRedis(userID uuid.UUID) (string, error) {
	ctx := context.Background()
	token, err := Rdb.Get(ctx, userID.String()).Result()
	if err != nil {
		return "", err
	}

	return token, nil
}
