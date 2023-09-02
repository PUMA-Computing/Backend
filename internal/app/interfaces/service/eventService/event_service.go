package eventService

import (
	"Backend/internal/app/domain/event"
	"Backend/internal/app/domain/user"
	event2 "Backend/internal/app/interfaces/repository/eventRepository"
	"errors"
	"github.com/google/uuid"
)

type EventService interface {
	CreateEvent(event *event.Events) error
	UpdateEvent(eventID int64, updatedEvent *event.Events) error
	DeleteEvent(eventID int64) error
	GetEventUsers(eventID int64) ([]*user.User, error)
	RegisterUserForEvent(UserID uuid.UUID, eventID int64) error
	GetEvent() ([]*event.Events, error)
	GetEventByID(eventID int64) (*event.Events, error)
}

type EventServiceImpl struct {
	eventRepository event2.EventRepository
}

func NewEventService(eventRepository event2.EventRepository) *EventServiceImpl {
	return &EventServiceImpl{eventRepository: eventRepository}
}

func (s *EventServiceImpl) CreateEvent(event *event.Events) error {
	return s.eventRepository.CreateEvent(event)
}

func (s *EventServiceImpl) UpdateEvent(eventID int64, updatedEvent *event.Events) error {
	existingEvent, err := s.eventRepository.GetEventByID(eventID)
	if err != nil {
		return err
	}

	existingEvent.Name = updatedEvent.Name
	existingEvent.Description = updatedEvent.Description
	existingEvent.Date = updatedEvent.Date

	return s.eventRepository.UpdateEvent(existingEvent)
}

func (s *EventServiceImpl) DeleteEvent(eventID int64) error {
	return s.eventRepository.DeleteEvent(eventID)
}

func (s *EventServiceImpl) GetEvent() ([]*event.Events, error) {
	events, err := s.eventRepository.GetEvent()
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *EventServiceImpl) GetEventByID(eventID int64) (*event.Events, error) {
	event, err := s.eventRepository.GetEventByID(eventID)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *EventServiceImpl) GetEventUsers(EventID int64) ([]*user.User, error) {
	users, err := s.eventRepository.GetEventUser(EventID)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *EventServiceImpl) RegisterUserForEvent(UserID uuid.UUID, eventID int64) error {
	event, err := s.eventRepository.GetEventByID(eventID)
	if err != nil {
		return err
	}

	if event == nil {
		return errors.New("event not found")
	}

	user, err := s.eventRepository.GetUserByID(UserID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	return s.eventRepository.RegisterUserForEvent(UserID, eventID)
}
