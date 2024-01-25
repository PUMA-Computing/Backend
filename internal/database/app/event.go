package app

import (
	"Backend/internal/database"
	"Backend/internal/models"
	"context"
	"github.com/google/uuid"
	"strconv"
	"time"
)

type Event struct {
}

func CreateEvent(event *models.Event) error {
	_, err := database.DB.Exec(context.Background(), `
        INSERT INTO events (title, description, start_date, end_date, user_id, status, link, thumbnail, category_id) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		event.Title, event.Description, event.StartDate, event.EndDate, event.UserID, event.Status, event.Link, event.Thumbnail, event.CategoryID)
	return err
}

func UpdateEvent(eventID int, updatedEvent *models.Event) error {
	_, err := database.DB.Exec(context.Background(), `
		UPDATE events SET title = $1, description = $2, start_date = $3, end_date = $4, user_id = $5, status = $6, link = $7, thumbnail = $8, category_id = $9
		WHERE id = $9`,
		updatedEvent.Title, updatedEvent.Description, updatedEvent.StartDate, updatedEvent.EndDate, updatedEvent.UserID, updatedEvent.Status, updatedEvent.Link, updatedEvent.Thumbnail, updatedEvent.CategoryID, eventID)
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
		SELECT id, title, description, start_date, end_date, user_id, status, link, thumbnail, created_at, updated_at, category_id
		FROM events WHERE id = $1`, eventID).Scan(&event.ID, &event.Title, &event.Description, &event.StartDate, &event.EndDate, &event.UserID, &event.Status, &event.Link, &event.Thumbnail, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func ListEvents(queryParams map[string]string) ([]*models.Event, error) {
	query := `
		SELECT id, title, description, start_date, end_date, user_id, status, link, thumbnail, created_at, updated_at, category_id
		FROM events
		WHERE TRUE`

	var args []interface{}

	if categoryID, ok := queryParams["category_id"]; ok {
		query += " AND category_id = $" + strconv.Itoa(len(args)+1)
		args = append(args, categoryID)
	}

	// Search by status
	if status, ok := queryParams["status"]; ok {
		query += " AND status = $" + strconv.Itoa(len(args)+1)
		args = append(args, status)
	}

	rows, err := database.DB.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID, &event.Title, &event.Description, &event.StartDate, &event.EndDate, &event.UserID, &event.Status, &event.Link, &event.Thumbnail, &event.CreatedAt, &event.UpdatedAt, &event.CategoryID)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}

func RegisterForEvent(userID uuid.UUID, eventID int) error {
	_, err := database.DB.Exec(context.Background(), `
		INSERT INTO event_registrations (event_id, user_id, registration_date) 
		VALUES ($1, $2, $3)`,
		eventID, userID, time.Now())
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
