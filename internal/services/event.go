package services

import (
	"Backend/internal/database/app"
	"Backend/internal/models"
	"github.com/google/uuid"
)

type EventService struct {
}

func NewEventService() *EventService {
	return &EventService{}
}

func (es *EventService) CreateEvent(event *models.Event) error {
	if err := app.CreateEvent(event); err != nil {
		return err
	}

	return nil
}

func (es *EventService) GetEventByID(eventID int) (*models.Event, error) {
	event, err := app.GetEventByID(eventID)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (es *EventService) EditEvent(eventID int, updatedEvent *models.Event) error {
	existingEvent, err := app.GetEventByID(eventID)
	if err != nil {
		return err
	}

	if updatedEvent.Title == "" {
		updatedEvent.Title = existingEvent.Title
	}

	if updatedEvent.Description == "" {
		updatedEvent.Description = existingEvent.Description
	}

	if updatedEvent.Date.IsZero() {
		updatedEvent.Date = existingEvent.Date
	}

	if updatedEvent.UserID == uuid.Nil {
		updatedEvent.UserID = existingEvent.UserID
	}

	if updatedEvent.CreatedAt.IsZero() {
		updatedEvent.CreatedAt = existingEvent.CreatedAt
	}

	if updatedEvent.UpdatedAt.IsZero() {
		updatedEvent.UpdatedAt = existingEvent.UpdatedAt
	}

	if err := app.UpdateEvent(eventID, updatedEvent); err != nil {
		return err
	}

	return nil
}

func (es *EventService) DeleteEvent(eventID int) error {
	if err := app.DeleteEvent(eventID); err != nil {
		return err
	}

	return nil
}

func (es *EventService) ListEvents() ([]*models.Event, error) {
	events, err := app.ListEvents()
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (es *EventService) RegisterForEvent(userID uuid.UUID, eventID int) error {
	if err := app.RegisterForEvent(userID, eventID); err != nil {
		return err
	}

	return nil
}

func (es *EventService) ListRegisteredUsers(eventID int) ([]*models.User, error) {
	users, err := app.ListRegisteredUsers(eventID)
	if err != nil {
		return nil, err
	}

	return users, nil
}
