package token

import (
	postgresRepository2 "Backend/internal/app/interfaces/repository/postgresRepository"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"strings"
	"time"
)

const (
	jwtExpirationDuration = time.Hour * 24
	SessionDuration       = 24 * time.Hour
)

var jwtSecretKey string

func GenerateJWTToken(userID uuid.UUID, roleID int) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID.String(),
		"exp":    time.Now().Add(jwtExpirationDuration).Unix(),
	}

	fmt.Printf("claims: %v\n", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecretKey))
}

func ExtractBearerToken(authHeader string) string {
	const bearerPrefix = "Bearer "
	if strings.HasPrefix(authHeader, bearerPrefix) {
		return strings.TrimPrefix(authHeader, bearerPrefix)
	}

	//if strings.HasPrefix(authHeader, bearerPrefix) {
	//	token := strings.TrimPrefix(authHeader, bearerPrefix)
	//	claims, err := ValidateSessionToken(token)
	//	if err != nil {
	//		fmt.Printf("Error extracting roleID from token: %v\n", err)
	//	} else {
	//		roleID, ok := claims["roleID"].(int)
	//		if !ok {
	//			fmt.Println("RoleID not found in token claims")
	//		} else {
	//			fmt.Printf("Extracted RoleID: %d\n", roleID)
	//		}
	//	}
	//	return token
	//}

	return ""
}

func IsValidSessionToken(sessionUserID, sessionToken string) (bool, error) {
	token, err := jwt.Parse(sessionToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, errors.New("invalid token claims")
	}

	userID, ok := claims["userID"].(string)
	if !ok || userID != sessionUserID {
		return false, errors.New("missing userID claim")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return false, errors.New("missing expiry claim")
	}

	expiryTime := time.Unix(int64(exp), 0)

	if time.Now().After(expiryTime) {
		return false, errors.New("token expired")
	}

	return true, nil
}

func IsSessionTokenAboutExpired(tokenString string, threshold time.Duration) bool {
	parsedToken, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
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

func ValidateSessionToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		return "", errors.New("missing userID claim")
	}

	fmt.Printf("ValidateSessionToken: userID: %s", userID)

	return userID, nil
}

func StoreSessionData(userID uuid.UUID, sessionToken string, expirationTime time.Time) error {
	postgresRepository, err := postgresRepository2.NewPostgresRepository()
	if err != nil {
		return err
	}
	defer postgresRepository.Close()

	if err := postgresRepository.StoreSessionData(userID, sessionToken, expirationTime); err != nil {
		return err
	}

	return nil
}

//func DeleteSessionData(userID uuid.UUID) error {
//	postgresRepository, err := postgresRepository2.NewPostgresRepository()
//	if err != nil {
//		return err
//	}
//	defer postgresRepository.Close()
//
//	if err := postgresRepository.DeleteSessionData(userID); err != nil {
//		return err
//	}
//
//	return nil
//}
