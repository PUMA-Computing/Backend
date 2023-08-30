package repository

import (
	"Backend/internal/app/domain"
	"github.com/gocql/gocql"
)

type EventRepository interface {
	IsRegisteredForEvent(userID, eventID string) (bool, error)
	GetEvent() ([]*domain.Event, error)
	GetEventByID(eventID string) (*domain.Event, error)
	UpdateEvent(event *domain.Event) error
	RegisterUserForEvent(userID, eventID string) error
}
type CassandraEventRepository struct {
	session *gocql.Session
}

func NewCassandraForEventRepository(session *gocql.Session) *CassandraEventRepository {
	return &CassandraEventRepository{session: session}
}

func (r *CassandraEventRepository) IsRegisteredForEvent(userID, eventID string) (bool, error) {
	query := r.session.Query(
		"SELECT COUNT(*) FROM user_id = ? AND event_id = ?",
		userID, eventID,
	)

	var count int
	if err := query.Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *CassandraEventRepository) GetEvent() ([]*domain.Event, error) {
	query := r.session.Query(
		"SELECT id, name, description, date, registered_users FROM events ",
	)
	iter := query.Iter()
	defer iter.Close()

	var events []*domain.Event

	for {
		var event domain.Event
		if !iter.Scan(&event.ID, &event.Name, &event.Description, &event.Date, &event.RegisteredUsers) {
			break
		}
		events = append(events, &event)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *CassandraEventRepository) GetEventByID(eventID string) (*domain.Event, error) {
	var event domain.Event

	query := r.session.Query(
		"SELECT id, name, description, date, registered_users FROM events WHERE id = ?",
		eventID,
	)

	if err := query.Scan(&event.ID, &event.Name, &event.Description, &event.Date, &event.RegisteredUsers); err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *CassandraEventRepository) UpdateEvent(event *domain.Event) error {
	query := r.session.Query(
		"UPDATE events SET name = ?, description = ?, date = ?, registered_users = ? WHERE id = ?",
		event.Name, event.Description, event.Date, event.RegisteredUsers, event.ID,
	)

	return query.Exec()
}

func (r *CassandraEventRepository) RegisterUserForEvent(userID, eventID string) error {
	query := r.session.Query(
		"INSERT INTO events (user_id, event_id) VALUES (?, ?)",
		userID, eventID,
	)

	return query.Exec()
}
