package utils

import (
	"math/rand"
	"time"
)

type StorageService int

const (
	AWSService StorageService = iota
	R2Service
)

// ChooseStorageService randomly selects a storage service based on a 75% chance of R2Service and 25% chance of AWSService
func ChooseStorageService() StorageService {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	serviceChoice := rand.Intn(100) + 1

	if serviceChoice <= 75 {
		return R2Service
	} else {
		return AWSService
	}
}
