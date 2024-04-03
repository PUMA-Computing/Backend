package utils

import (
	"crypto/rand"
)

type StorageService int

const (
	AWSService StorageService = iota
	R2Service
)

// ChooseStorageService randomly selects a storage service based on a 75% chance of R2Service and 25% chance of AWSService
func ChooseStorageService() StorageService {
	// Generate a random number between 0 and 3
	randomBytes := make([]byte, 1)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	randomNumber := randomBytes[0] % 4

	if randomNumber == 0 {
		return AWSService
	}
	return R2Service
}
