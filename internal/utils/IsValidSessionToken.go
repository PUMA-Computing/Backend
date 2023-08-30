package utils

import (
	"Backend/internal/app/repository"
	"github.com/gocql/gocql"
	"time"
)

func IsValidSessionToken(sessionUserID, sessionToken string) (bool, error) {
	userID, _, err := ValidateSessionToken(sessionToken)
	if err != nil {
		return false, err
	}

	userUUID, err := gocql.ParseUUID(userID)
	if err != nil {
		return false, err
	}

	repo, err := repository.NewCassandraRepository()
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
