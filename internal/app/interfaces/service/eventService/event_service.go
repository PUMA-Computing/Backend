package eventService

import (
	"Backend/internal/app/domain/event"
	"Backend/internal/app/domain/user"
	event2 "Backend/internal/app/interfaces/repository/eventRepository"
	"errors"
)

type EventService interface {
	CreateEvent(event *event.Event) error
	UpdateEvent(eventID string, updatedEvent *event.Event) error
	DeleteEvent(eventID string) error
	GetEventUsers(eventID string) ([]*user.User, error)
	RegisterUserForEvent(UserID, eventID string) error
	GetEvent() ([]*event.Event, error)
	GetEventByID(eventID string) (*event.Event, error)
}

type EventServiceImpl struct {
	eventRepository event2.EventRepository
}

func NewEventService(eventRepository event2.EventRepository) *EventServiceImpl {
	return &EventServiceImpl{eventRepository: eventRepository}
}

func (s *EventServiceImpl) CreateEvent(event *event.Event) error {
	return s.eventRepository.CreateEvent(event)
}

func (s *EventServiceImpl) UpdateEvent(eventID string, updatedEvent *event.Event) error {
	existingEvent, err := s.eventRepository.GetEventByID(eventID)
	if err != nil {
		return err
	}

	existingEvent.Name = updatedEvent.Name
	existingEvent.Description = updatedEvent.Description
	existingEvent.Date = updatedEvent.Date

	return s.eventRepository.UpdateEvent(existingEvent)
}

func (s *EventServiceImpl) DeleteEvent(eventID string) error {
	return s.eventRepository.DeleteEvent(eventID)
}

func (s *EventServiceImpl) GetEvent() ([]*event.Event, error) {
	events, err := s.eventRepository.GetEvent()
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *EventServiceImpl) GetEventByID(eventID string) (*event.Event, error) {
	event, err := s.eventRepository.GetEventByID(eventID)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *EventServiceImpl) GetEventUsers(EventID string) ([]*user.User, error) {
	users, err := s.eventRepository.GetEventUser(EventID)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *EventServiceImpl) RegisterUserForEvent(UserID, eventID string) error {
	isRegistered, err := s.eventRepository.IsRegisteredForEvent(UserID, eventID)
	if err != nil {
		return err
	}

	if isRegistered {
		return errors.New("user already registered for event")
	}

	event, err := s.eventRepository.GetEventByID(eventID)
	if err != nil {
		return err
	}

	event.RegisteredUsers = append(event.RegisteredUsers, UserID)
	if err := s.eventRepository.UpdateEvent(event); err != nil {
		return err
	}

	return nil
}
