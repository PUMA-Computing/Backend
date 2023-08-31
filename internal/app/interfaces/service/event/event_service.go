package event

import (
	"Backend/internal/app/domain/event"
	"Backend/internal/app/domain/user"
	event2 "Backend/internal/app/interfaces/repository/event"
	"errors"
)

type EventsService interface {
	RegisterUserForEvent(UserID, eventID string) error
	GetEvent() ([]*event.Event, error)
	GetEventByID(eventID string) (*event.Event, error)
}

type EventService struct {
	eventRepository event2.EventRepository
}

func NewEventService(eventRepository event2.EventRepository) *EventService {
	return &EventService{eventRepository: eventRepository}
}

func (s *EventService) CreateEvent(event *event.Event) error {
	return s.eventRepository.CreateEvent(event)
}

func (s *EventService) UpdateEvent(eventID string, updatedEvent *event.Event) error {
	existingEvent, err := s.eventRepository.GetEventByID(eventID)
	if err != nil {
		return err
	}

	existingEvent.Name = updatedEvent.Name
	existingEvent.Description = updatedEvent.Description
	existingEvent.Date = updatedEvent.Date

	return s.eventRepository.UpdateEvent(existingEvent)
}

func (s *EventService) DeleteEvent(eventID string) error {
	return s.eventRepository.DeleteEvent(eventID)
}

func (s *EventService) GetEvent() ([]*event.Event, error) {
	events, err := s.eventRepository.GetEvent()
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *EventService) GetEventUsers(EventID string) ([]*user.User, error) {
	users, err := s.eventRepository.GetEventUser(EventID)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *EventService) RegisterUserForEvent(UserID, eventID string) error {
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
