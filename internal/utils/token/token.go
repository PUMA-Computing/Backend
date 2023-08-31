package token

import (
	"Backend/internal/app/interfaces/repository/cassandraRepository"
	"errors"
	"github.com/gocql/gocql"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	jwtSecretKey          = "pumacomputing"
	jwtExpirationDuration = time.Hour * 24
	SessionDuration       = 24 * time.Hour
)

func GenerateJWTToken(userID gocql.UUID, role string) (string, error) {
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

func IsValidSessionToken(sessionUserID, sessionToken string) (bool, error) {
	userID, _, err := ValidateSessionToken(sessionToken)
	if err != nil {
		return false, err
	}

	userUUID, err := gocql.ParseUUID(userID)
	if err != nil {
		return false, err
	}

	repo, err := cassandraRepository.NewCassandraRepository()
	if err != nil {
		return false, err
	}
	defer repo.Close()

	sessionData, err := repo.GetSessionData(userUUID)
	if err != nil {
		return false, err
	}

	if sessionData.SessionToken != sessionToken {
		return false, nil
	}

	if time.Now().After(sessionData.ExpiredAt) {
		return false, nil
	}

	return true, nil
}

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

func StoreSessionData(userID gocql.UUID, sessionToken string, expirationTime time.Time) error {
	cassandraRepository, err := cassandraRepository.NewCassandraRepository()
	if err != nil {
		return err
	}
	defer cassandraRepository.Close()

	if err := cassandraRepository.StoreSessionData(userID, sessionToken, expirationTime); err != nil {
		return err
	}

	return nil
}
