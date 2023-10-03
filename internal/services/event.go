package services

import (
	"Backend/internal/database"
	"Backend/internal/models"
	"github.com/google/uuid"
)

type EventService struct {
}

func NewEventService() *EventService {
	return &EventService{}
}

func (es *EventService) CreateEvent(event *models.Event) error {
	if err := database.CreateEvent(event); err != nil {
		return err
	}

	return nil
}

func (s *EventService) GetEventByID(eventID int) (*models.Event, error) {
	event, err := database.GetEventByID(eventID)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (es *EventService) EditEvent(eventID int, updatedEvent *models.Event) error {
	if err := database.UpdateEvent(eventID, updatedEvent); err != nil {
		return err
	}

	return nil
}

func (es *EventService) DeleteEvent(eventID int) error {
	if err := database.DeleteEvent(eventID); err != nil {
		return err
	}

	return nil
}

func (es *EventService) ListEvents() ([]*models.Event, error) {
	events, err := database.ListEvents()
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (es *EventService) RegisterForEvent(userID uuid.UUID, eventID int) error {
	if err := database.RegisterForEvent(userID, eventID); err != nil {
		return err
	}

	return nil
}

func (es *EventService) ListRegisteredUsers(eventID int) ([]*models.User, error) {
	users, err := database.ListRegisteredUsers(eventID)
	if err != nil {
		return nil, err
	}

	return users, nil
}
