package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func IsTokenAboutToExpire(token string, threshold time.Duration) bool {
	parsedToken, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return false
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	expiry, ok := claims["exp"].(float64)
	if !ok {
		return false
	}

	expiryTime := time.Unix(int64(expiry), 0)
	currentTime := time.Now()

	return currentTime.Add(threshold).After(expiryTime)
}
