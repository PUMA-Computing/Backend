package eventRepository

import (
	"Backend/internal/app/domain/event"
	"Backend/internal/app/domain/user"
	"github.com/gocql/gocql"
)

type EventRepository interface {
	IsRegisteredForEvent(userID gocql.UUID, eventID int64) (bool, error)
	GetEvent() ([]*event.Event, error)
	GetEventByID(eventID int64) (*event.Event, error)
	GetEventUser(eventID int64) ([]*user.User, error)
	GetUserByID(userID gocql.UUID) (*user.User, error)
	CreateEvent(event *event.Event) error
	UpdateEvent(event *event.Event) error
	DeleteEvent(eventID int64) error
	RegisterUserForEvent(userID gocql.UUID, eventID int64) error
}
type CassandraEventRepository struct {
	session *gocql.Session
}

func NewCassandraForEventRepository(session *gocql.Session) *CassandraEventRepository {
	return &CassandraEventRepository{session: session}
}

func (r *CassandraEventRepository) IsRegisteredForEvent(userID gocql.UUID, eventID int64) (bool, error) {
	query := r.session.Query(
		"SELECT COUNT(*) FROM events WHERE registered_users = ? AND id = ?",
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
		"SELECT id, image, name, description, date, registered_users FROM events ",
	)
	iter := query.Iter()
	defer iter.Close()

	var events []*event.Event

	for {
		var event event.Event
		if !iter.Scan(&event.ID, &event.Image, &event.Name, &event.Description, &event.Date, &event.RegisteredUsers) {
			break
		}
		events = append(events, &event)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *CassandraEventRepository) GetEventByID(eventID int64) (*event.Event, error) {
	var event event.Event

	query := r.session.Query(
		"SELECT id, image, name, description, date, registered_users FROM events WHERE id = ?",
		eventID,
	)

	if err := query.Scan(&event.ID, &event.Image, &event.Name, &event.Description, &event.Date, &event.RegisteredUsers); err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *CassandraEventRepository) GetEventUser(eventID int64) ([]*user.User, error) {
	query := r.session.Query(
		"SELECT id, first_name, last_name, email, nim, year, role FROM users WHERE id = ?",
		eventID,
	)

	var registeredUsers []gocql.UUID
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

func (r *CassandraEventRepository) GetUserByID(userID gocql.UUID) (*user.User, error) {
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
		"INSERT INTO events (id, image, name, description, date, registered_users) VALUES (?, ?, ?, ?, ?)",
		event.ID, event.Image, event.Name, event.Description, event.Date, event.RegisteredUsers,
	)

	return query.Exec()
}

func (r *CassandraEventRepository) UpdateEvent(event *event.Event) error {
	query := r.session.Query(
		"UPDATE events SET image = ?, name = ?, description = ?, date = ?, registered_users = ? WHERE id = ?",
		event.Image, event.Name, event.Description, event.Date, event.RegisteredUsers, event.ID,
	)

	return query.Exec()
}

func (r *CassandraEventRepository) DeleteEvent(eventID int64) error {
	query := r.session.Query(
		"DELETE FROM events WHERE id = ?",
		eventID,
	)

	return query.Exec()
}

func (r *CassandraEventRepository) RegisterUserForEvent(userID gocql.UUID, eventID int64) error {
	query := r.session.Query(
		"INSERT INTO events (registered_users, id) VALUES (?, ?)",
		userID, eventID,
	)

	return query.Exec()
}
