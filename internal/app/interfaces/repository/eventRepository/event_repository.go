package eventRepository

import (
	"Backend/internal/app/domain/event"
	"Backend/internal/app/domain/user"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type EventRepository interface {
	IsRegisteredForEvent(userID uuid.UUID, eventID int64) (bool, error)
	GetEvent() ([]*event.Events, error)
	GetEventByID(eventID int64) (*event.Events, error)
	GetEventUser(eventID int64) ([]*user.User, error)
	GetUserByID(userID uuid.UUID) (*user.User, error)
	CreateEvent(event *event.Events) error
	UpdateEvent(event *event.Events) error
	DeleteEvent(eventID int64) error
	RegisterUserForEvent(userID uuid.UUID, eventID int64) error
}
type CassandraEventRepository struct {
	session *gocql.Session
}

type PostgresEventRepository struct {
	DB *gorm.DB
}

func NewPostgresForEventRepository(DB *gorm.DB) *PostgresEventRepository {
	return &PostgresEventRepository{DB: DB}
}

func (r *PostgresEventRepository) GetEvent() ([]*event.Events, error) {
	var events []*event.Events
	if err := r.DB.Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

func (r *PostgresEventRepository) GetEventByID(eventID int64) (*event.Events, error) {
	var event event.Events
	if err := r.DB.Where("id = ?", eventID).First(&event).Error; err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *PostgresEventRepository) GetEventUser(eventID int64) ([]*user.User, error) {
	var users []*user.User
	if err := r.DB.
		Joins("JOIN event_registration ON event_registration.user_id = users.id").
		Where("event_registration.event_id = ?", eventID).
		Find(&users).
		Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *PostgresEventRepository) GetUserByID(userID uuid.UUID) (*user.User, error) {
	var u user.User
	if err := r.DB.Where("id = ?", userID).First(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *PostgresEventRepository) RegisterUserForEvent(userID uuid.UUID, eventID int64) error {
	registration := event.Registration{
		EventID:   eventID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	return r.DB.Create(&registration).Error
}

func (r *PostgresEventRepository) IsRegisteredForEvent(userID uuid.UUID, eventID int64) (bool, error) {
	var count int64
	if err := r.DB.
		Model(&event.Registration{}).
		Where("user_id = ? AND event_id = ?", userID, eventID).
		Count(&count).
		Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PostgresEventRepository) CreateEvent(event *event.Events) error {
	return r.DB.Create(event).Error
}

func (r *PostgresEventRepository) UpdateEvent(event *event.Events) error {
	return r.DB.Save(event).Error
}

func (r *PostgresEventRepository) DeleteEvent(eventID int64) error {
	return r.DB.Delete(&event.Events{}, eventID).Error
}
