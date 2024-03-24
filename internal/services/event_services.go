package services

import (
	"Backend/internal/database/app"
	"Backend/internal/models"
	"github.com/google/uuid"
	"time"
)

type EventService struct {
}

func NewEventService() *EventService {
	return &EventService{}
}

func (es *EventService) CreateEvent(event *models.Event) error {
	if time.Now().Before(event.StartDate) {
		event.Status = "Upcoming"
	} else if time.Now().After(event.StartDate) && time.Now().Before(event.EndDate) {
		event.Status = "Open"
	} else {
		event.Status = "Ended"
	}

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

func (es *EventService) GetEventBySlug(slug string) (*models.Event, error) {
	event, err := app.GetEventBySlug(slug)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (es *EventService) EditEvent(eventID int, updatedEvent *models.Event) error {
	if time.Now().Before(updatedEvent.StartDate) {
		updatedEvent.Status = "Upcoming"
	} else if time.Now().After(updatedEvent.StartDate) && time.Now().Before(updatedEvent.EndDate) {
		updatedEvent.Status = "Open"
	} else {
		updatedEvent.Status = "Ended"
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

func (es *EventService) ListEvents(queryParams map[string]string) ([]*models.Event, error) {
	events, err := app.ListEvents(queryParams)
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

func (es *EventService) ListEventsRegisteredByUser(userID uuid.UUID) ([]*models.Event, error) {
	events, err := app.ListEventsRegisteredByUser(userID)
	if err != nil {
		return nil, err
	}

	return events, nil
}
