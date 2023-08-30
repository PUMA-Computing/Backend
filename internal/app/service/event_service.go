package service

import (
	"Backend/internal/app/domain"
	"Backend/internal/app/repository"
	"errors"
)

type EventsService interface {
	RegisterUserForEvent(UserID, eventID string) error
	GetEvent() ([]*domain.Event, error)
	GetEventByID(eventID string) (*domain.Event, error)
}

type EventService struct {
	eventRepository repository.EventRepository
}

func NewEventService(eventRepository repository.EventRepository) *EventService {
	return &EventService{eventRepository: eventRepository}
}

func (s *EventService) GetEvent() ([]*domain.Event, error) {
	events, err := s.eventRepository.GetEvent()
	if err != nil {
		return nil, err
	}

	return events, nil
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
