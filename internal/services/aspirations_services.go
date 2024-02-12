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

func (s *AspirationService) CloseAspirationByID(id int) error {
	return app.CloseAspirationByID(id)
}

func (s *AspirationService) DeleteAspirationByID(id int) error {
	return app.DeleteAspirationByID(id)
}

func (s *AspirationService) GetAspirations(queryParams map[string]string) ([]models.Aspiration, error) {
	return app.GetAspirations(queryParams)
}

func (s *AspirationService) GetAspirationByID(id int) (*models.Aspiration, error) {
	return app.GetAspirationByID(id)
}
