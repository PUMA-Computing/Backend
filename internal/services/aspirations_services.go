package services

import (
	"Backend/internal/database/app"
	"Backend/internal/models"
)

type AspirationService struct{}

func NewAspirationService() *AspirationService {
	return &AspirationService{}
}

func (s *AspirationService) CreateAspiration(aspiration *models.Aspiration) error {
	return app.CreateAspiration(aspiration)
}
