package services

import (
	"Backend/internal/database/app"
	"Backend/internal/models"
	"github.com/google/uuid"
)

type AspirationService struct{}

func NewAspirationService() *AspirationService {
	return &AspirationService{}
}

func (s *AspirationService) CreateAspiration(aspiration *models.Aspiration) (*models.Aspiration, error) {
	createdAspiration, err := app.CreateAspiration(aspiration)
	if err != nil {
		return nil, err
	}

	// Add upvote to the aspiration by the user who created it
	err = app.AddUpvote(aspiration.UserID, createdAspiration.ID)
	if err != nil {
		return nil, err
	}

	return createdAspiration, nil
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

func (s *AspirationService) UpvoteExists(userID uuid.UUID, aspirationID int) (bool, error) {
	return app.UpvoteExists(userID, aspirationID)
}

func (s *AspirationService) AddUpvote(userID uuid.UUID, aspirationID int) error {
	return app.AddUpvote(userID, aspirationID)
}

func (s *AspirationService) RemoveUpvote(userID uuid.UUID, aspirationID int) error {
	return app.RemoveUpvote(userID, aspirationID)
}

func (s *AspirationService) GetUpvotesByAspirationID(aspirationID int) (int, error) {
	return app.GetUpvotesByAspirationID(aspirationID)
}

func (s *AspirationService) AddAdminReply(aspirationID int, adminReply string) error {
	// Close the aspiration
	err := s.CloseAspirationByID(aspirationID)
	if err != nil {
		return err
	}
	return app.AddAdminReply(aspirationID, adminReply)
}
