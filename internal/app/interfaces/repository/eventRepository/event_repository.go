package eventRepository

import (
	"Backend/internal/app/domain/event"
	"Backend/internal/app/domain/user"
	"github.com/gocql/gocql"
)

type EventRepository interface {
	IsRegisteredForEvent(userID, eventID string) (bool, error)
	GetEvent() ([]*event.Event, error)
	GetEventByID(eventID string) (*event.Event, error)
	GetEventUser(eventID string) ([]*user.User, error)
	GetUserByID(userID string) (*user.User, error)
	CreateEvent(event *event.Event) error
	UpdateEvent(event *event.Event) error
	DeleteEvent(eventID string) error
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

func (r *CassandraEventRepository) GetEvent() ([]*event.Event, error) {
	query := r.session.Query(
		"SELECT id, name, description, date, registered_users FROM events ",
	)
	iter := query.Iter()
	defer iter.Close()

	var events []*event.Event

	for {
		var event event.Event
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

func (r *CassandraEventRepository) GetEventByID(eventID string) (*event.Event, error) {
	var event event.Event

	query := r.session.Query(
		"SELECT id, name, description, date, registered_users FROM events WHERE id = ?",
		eventID,
	)

	if err := query.Scan(&event.ID, &event.Name, &event.Description, &event.Date, &event.RegisteredUsers); err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *CassandraEventRepository) GetEventUser(eventID string) ([]*user.User, error) {
	query := r.session.Query(
		"SELECT id, first_name, last_name, email, nim, year, role FROM users WHERE id = ?",
		eventID,
	)

	var registeredUsers []string
	if err := query.Scan(&registeredUsers); err != nil {
		return nil, err
	}

	users := make([]*user.User, 0)
	for _, UserID := range registeredUsers {
		user, err := r.GetUserByID(UserID)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}

func (r *CassandraEventRepository) GetUserByID(userID string) (*user.User, error) {
	var user user.User

	query := r.session.Query(
		"SELECT id, first_name, last_name, nim, email, year, role FROM users WHERE id = ?",
		userID,
	)

	if err := query.Scan(&userID, &user.FirstName, &user.LastName, &user.NIM, &user.Email, &user.Year, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil

}

func (r *CassandraEventRepository) CreateEvent(event *event.Event) error {
	query := r.session.Query(
		"INSERT INTO events (id, name, description, date, registered_users) VALUES (?, ?, ?, ?, ?)",
		event.ID, event.Name, event.Description, event.Date, event.RegisteredUsers,
	)

	return query.Exec()
}

func (r *CassandraEventRepository) UpdateEvent(event *event.Event) error {
	query := r.session.Query(
		"UPDATE events SET name = ?, description = ?, date = ?, registered_users = ? WHERE id = ?",
		event.Name, event.Description, event.Date, event.RegisteredUsers, event.ID,
	)

	return query.Exec()
}

func (r *CassandraEventRepository) DeleteEvent(eventID string) error {
	query := r.session.Query(
		"DELETE FROM events WHERE id = ?",
		eventID,
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
