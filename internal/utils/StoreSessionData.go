package utils

import (
	"Backend/internal/app/repository"
	"github.com/gocql/gocql"
	"time"
)

func StoreSessionData(userID gocql.UUID, sessionToken string, expirationTime time.Time) error {
	cassandraRepository, err := repository.NewCassandraRepository()
	if err != nil {
		return err
	}
	defer cassandraRepository.Close()

	if err := cassandraRepository.StoreSessionData(userID, sessionToken, expirationTime); err != nil {
		return err
	}

	return nil
}
