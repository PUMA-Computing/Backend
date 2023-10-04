package app

import (
	"Backend/internal/database"
	"Backend/internal/models"
	"context"
	"github.com/google/uuid"
)

type Event struct {
}

func CreateEvent(event *models.Event) error {
	_, err := database.DB.Exec(context.Background(), `
        INSERT INTO events (title, description, date, user_id) 
        VALUES ($1, $2, $3, $4)`,
		event.Title, event.Description, event.Date, event.UserID)
	return err
}

func UpdateEvent(eventID int, updatedEvent *models.Event) error {
	_, err := database.DB.Exec(context.Background(), `
		UPDATE events SET title = $1, description = $2, date = $3, user_id = $4
		WHERE id = $5`,
		updatedEvent.Title, updatedEvent.Description, updatedEvent.Date, updatedEvent.UserID, eventID)
	return err
}

func DeleteEvent(eventID int) error {
	_, err := database.DB.Exec(context.Background(), `
		DELETE FROM events WHERE id = $1`, eventID)
	return err
}

func GetEventByID(eventID int) (*models.Event, error) {
	var event models.Event
	err := database.DB.QueryRow(context.Background(), `
		SELECT id, title, description, date, user_id, created_at, updated_at
		FROM events WHERE id = $1`, eventID).Scan(&event.ID, &event.Title, &event.Description, &event.Date, &event.UserID, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func ListEvents() ([]*models.Event, error) {
	var events []*models.Event
	rows, err := database.DB.Query(context.Background(), `
		SELECT id, title, description, date, user_id, created_at, updated_at
		FROM events`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.Date, &event.UserID, &event.CreatedAt, &event.UpdatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}

func RegisterForEvent(userID uuid.UUID, eventID int) error {
	_, err := database.DB.Exec(context.Background(), `
		INSERT INTO event_registrations (event_id, user_id) 
		VALUES ($1, $2)`,
		eventID, userID)
	return err
}

func ListRegisteredUsers(eventID int) ([]*models.User, error) {
	rows, err := database.DB.Query(context.Background(), `
        SELECT u.id, u.username, u.first_name, u.middle_name, u.last_name, u.email, u.role_id, u.created_at, u.updated_at
        FROM users u
        JOIN event_registrations er ON u.id = er.user_id
        WHERE er.event_id = $1`,
		eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.RoleID,
			&user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}
