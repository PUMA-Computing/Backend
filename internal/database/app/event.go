package app

import (
	"Backend/internal/database"
	"Backend/internal/models"
	"Backend/pkg/utils"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"log"
	"time"
)

type Event struct {
}

func CreateEvent(event *models.Event) error {
	_, err := database.DB.Exec(context.Background(), `
        INSERT INTO events (title, description, start_date, end_date, user_id, status, slug, thumbnail, organization_id, max_registration) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		event.Title, event.Description, event.StartDate, event.EndDate, event.UserID, event.Status, event.Slug, event.Thumbnail, event.OrganizationID, event.MaxRegistration)
	return err
}

func UpdateEvent(eventID int, updatedEvent *models.Event) error {
	query := `
        UPDATE events SET title = $1, description = $2, start_date = $3, end_date = $4, user_id = $5, status = $6, slug = $7, thumbnail = $8, organization_id = $9, max_registration = $10, updated_at = $11 WHERE id = $12`

	// Log the SQL query and parameters
	log.Printf("Executing query: %s\n", query)
	log.Printf("Parameters: %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v\n",
		updatedEvent.Title, updatedEvent.Description, updatedEvent.StartDate,
		updatedEvent.EndDate, updatedEvent.UserID, updatedEvent.Status, updatedEvent.Slug,
		updatedEvent.Thumbnail, updatedEvent.OrganizationID, updatedEvent.MaxRegistration, time.Now(), eventID)

	_, err := database.DB.Exec(context.Background(), query,
		updatedEvent.Title, updatedEvent.Description, updatedEvent.StartDate,
		updatedEvent.EndDate, updatedEvent.UserID, updatedEvent.Status, updatedEvent.Slug,
		updatedEvent.Thumbnail, updatedEvent.OrganizationID, updatedEvent.MaxRegistration, time.Now(), eventID)

	if err != nil {
		log.Printf("Error executing query: %s\n", err.Error())
	}

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
		SELECT id, title, description, start_date, end_date, user_id, status, slug, thumbnail, created_at, updated_at, organization_id, max_registration
		FROM events WHERE id = $1`, eventID).Scan(&event.ID, &event.Title, &event.Description, &event.StartDate, &event.EndDate, &event.UserID, &event.Status, &event.Slug, &event.Thumbnail, &event.CreatedAt, &event.UpdatedAt, &event.OrganizationID, &event.MaxRegistration)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func GetEventBySlug(slug string) (*models.Event, error) {
	var event models.Event
	err := database.DB.QueryRow(context.Background(), `
		SELECT e.id, e.title, e.description, e.start_date, e.end_date, e.user_id, e.status, e.slug, e.thumbnail, e.created_at, e.updated_at, e.organization_id, e.max_registration, o.name AS organization, CONCAT(u.first_name, ' ', u.last_name) AS author
		FROM events e
		LEFT JOIN organizations o ON e.organization_id = o.id
		LEFT JOIN users u ON e.user_id = u.id
		WHERE e.slug = $1`, slug).Scan(&event.ID, &event.Title, &event.Description, &event.StartDate, &event.EndDate, &event.UserID, &event.Status, &event.Slug, &event.Thumbnail, &event.CreatedAt, &event.UpdatedAt, &event.OrganizationID, &event.MaxRegistration, &event.Organization, &event.Author)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// ListEvents returns a list of events based on the query parameters
func ListEvents(queryParams map[string]string) ([]*models.Event, error) {
	query := `
		SELECT events.*, organizations.name AS organization, CONCAT(users.first_name, ' ', users.last_name) AS author
		FROM events
		LEFT JOIN organizations ON events.organization_id = organizations.id
		LEFT JOIN users ON events.user_id = users.id
		WHERE 1=1`

	if queryParams["organization_id"] != "" {
		query += " AND events.organization_id = '" + queryParams["organization_id"] + "'"
	}

	// Search by status
	if queryParams["status"] != "" {
		query += " AND events.status = '" + queryParams["status"] + "'"
	}

	// Search by Slug
	if queryParams["slug"] != "" {
		query += " AND events.slug = '" + queryParams["slug"] + "'"
	}

	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartDate,
			&event.EndDate,
			&event.UserID,
			&event.Status,
			&event.Slug,
			&event.Thumbnail,
			&event.CreatedAt,
			&event.UpdatedAt,
			&event.OrganizationID,
			&event.MaxRegistration,
			&event.Organization,
			&event.Author)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}

// RegisterForEvent registers a user for an event by creating a new event registration record
func RegisterForEvent(userID uuid.UUID, eventID int) error {
	tx, err := database.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			err := tx.Rollback(context.Background())
			if err != nil {
				return
			}
			panic(p)
		} else if err != nil {
			err := tx.Rollback(context.Background())
			if err != nil {
				return
			}
		} else {
			err := tx.Commit(context.Background())
			if err != nil {
				return
			}
		}
	}()

	// Check if the event has a maximum registration limit
	var maxRegistration *int
	err = tx.QueryRow(context.Background(), `
        SELECT max_registration FROM events WHERE id = $1`, eventID).Scan(&maxRegistration)
	if err != nil {
		// Check if the error is due to no rows being returned
		if errors.Is(err, sql.ErrNoRows) {
			// No registration limit specified for the event, proceed with registration
			_, err = tx.Exec(context.Background(), `
                INSERT INTO event_registrations (event_id, user_id, registration_date)
                VALUES ($1, $2, $3)`, eventID, userID, time.Now())
			return err
		}
		return err
	}

	if maxRegistration != nil && *maxRegistration > 0 {
		// Check if the maximum registration limit has been reached
		var count int
		err := database.DB.QueryRow(context.Background(), `
            SELECT COUNT(*) FROM event_registrations WHERE event_id = $1`, eventID).Scan(&count)
		if err != nil {
			return err
		}

		if count >= *maxRegistration {
			return utils.MaxRegistrationReachedError{EventID: eventID}
		}
	}

	_, err = database.DB.Exec(context.Background(), `
        INSERT INTO event_registrations (event_id, user_id, registration_date)
        VALUES ($1, $2, $3)`, eventID, userID, time.Now())
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

func ListEventsRegisteredByUser(userID uuid.UUID) ([]*models.Event, error) {
	rows, err := database.DB.Query(context.Background(), `
		SELECT e.id, e.title, e.description, e.start_date, e.end_date, e.user_id, e.status, e.slug, e.thumbnail, e.created_at, e.updated_at, e.organization_id, e.max_registration
		FROM events e
		JOIN event_registrations er ON e.id = er.event_id
		WHERE er.user_id = $1`,
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID, &event.Title, &event.Description, &event.StartDate, &event.EndDate, &event.UserID, &event.Status, &event.Slug, &event.Thumbnail, &event.CreatedAt, &event.UpdatedAt, &event.OrganizationID, &event.MaxRegistration)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}
