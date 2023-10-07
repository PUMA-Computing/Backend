package services

import (
	"Backend/internal/database/app"
	"Backend/internal/models"
	"Backend/pkg/utils"
	"github.com/google/uuid"
	"log"
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

	if updatedEvent.Title != "" && updatedEvent.Title != existingEvent.Title {
		updatedEvent.Link = "/events/" + utils.GenerateFriendlyURL(updatedEvent.Title)
	} else {
		updatedEvent.Link = existingEvent.Link
	}

	if updatedEvent.Description == "" {
		updatedEvent.Description = existingEvent.Description
	}

	if updatedEvent.StartDate.IsZero() {
		updatedEvent.StartDate = existingEvent.StartDate
	}

	if updatedEvent.EndDate.IsZero() {
		updatedEvent.EndDate = existingEvent.EndDate
	}

	if updatedEvent.UserID == uuid.Nil {
		updatedEvent.UserID = existingEvent.UserID
	}

	if updatedEvent.Status == "" {
		updatedEvent.Status = existingEvent.Status
	}

	if updatedEvent.Link == "" {
		updatedEvent.Link = existingEvent.Link
	}

	if updatedEvent.CreatedAt.IsZero() {
		updatedEvent.CreatedAt = existingEvent.CreatedAt
	}

	if updatedEvent.UpdatedAt.IsZero() {
		updatedEvent.UpdatedAt = existingEvent.UpdatedAt
	}

	if updatedEvent.Thumbnail == "" {
		updatedEvent.Thumbnail = existingEvent.Thumbnail
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
	log.Println("Service ListEvents Begin")
	events, err := app.ListEvents()
	if err != nil {
		return nil, err
	}

	log.Println("Service ListEvents End")

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
