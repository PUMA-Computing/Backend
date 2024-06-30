package app

import (
	"Backend/internal/database"
	"Backend/internal/models"
	"Backend/pkg/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strconv"
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

// UpdateEvent updates an existing event record in the database with partial data
func UpdateEvent(eventID int, updatedEvent *models.Event) error {
	// Start building the query
	query := "UPDATE events SET "
	params := []interface{}{}
	paramCount := 1

	// Check each field in updatedEvent. If it's not the zero value for its type, add it to the query.
	if updatedEvent.Title != "" {
		query += "title = $" + strconv.Itoa(paramCount) + ", "
		params = append(params, updatedEvent.Title)
		paramCount++
	}
	if updatedEvent.Description != "" {
		query += "description = $" + strconv.Itoa(paramCount) + ", "
		params = append(params, updatedEvent.Description)
		paramCount++
	}
	if !updatedEvent.StartDate.IsZero() {
		query += "start_date = $" + strconv.Itoa(paramCount) + ", "
		params = append(params, updatedEvent.StartDate)
		paramCount++
	}
	if !updatedEvent.EndDate.IsZero() {
		query += "end_date = $" + strconv.Itoa(paramCount) + ", "
		params = append(params, updatedEvent.EndDate)
		paramCount++
	}
	if updatedEvent.Status != "" {
		query += "status = $" + strconv.Itoa(paramCount) + ", "
		params = append(params, updatedEvent.Status)
		paramCount++
	}
	if updatedEvent.Slug != "" {
		query += "slug = $" + strconv.Itoa(paramCount) + ", "
		params = append(params, updatedEvent.Slug)
		paramCount++
	}
	if updatedEvent.Thumbnail != "" {
		query += "thumbnail = $" + strconv.Itoa(paramCount) + ", "
		params = append(params, updatedEvent.Thumbnail)
		paramCount++
	}
	if updatedEvent.OrganizationID != 0 {
		query += "organization_id = $" + strconv.Itoa(paramCount) + ", "
		params = append(params, updatedEvent.OrganizationID)
		paramCount++
	}
	if updatedEvent.MaxRegistration != nil {
		query += "max_registration = $" + strconv.Itoa(paramCount) + ", "
		params = append(params, updatedEvent.MaxRegistration)
		paramCount++
	}

	// Remove the trailing comma and space, then add the WHERE clause
	query = query[:len(query)-2] + " WHERE id = $" + strconv.Itoa(paramCount)
	params = append(params, eventID)

	_, err := database.DB.Exec(context.Background(), query, params...)
	return err
}

// DeleteEvent deletes an event record from the database and delete all event registrations associated with the event
func DeleteEvent(eventID int) error {
	_, err := database.DB.Exec(context.Background(), `
		DELETE FROM event_registrations WHERE event_id = $1`, eventID)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec(context.Background(), `
		DELETE FROM events WHERE id = $1`, eventID)
	return err
}

// GetEventByID Get event and join the total number of registered users
func GetEventByID(eventID int) (*models.Event, error) {
	var event models.Event
	err := database.DB.QueryRow(context.Background(), `
		SELECT e.id, e.title, e.description, e.start_date, e.end_date, e.user_id, e.status, e.slug, e.thumbnail, e.created_at, e.updated_at, e.organization_id, e.max_registration, o.name as organization, CONCAT(u.first_name, ' ', u.last_name) AS author, COUNT(er.user_id) as total_registered
		FROM events e
		LEFT JOIN organizations o ON e.organization_id = o.id
		LEFT JOIN users u ON e.user_id = u.id
		LEFT JOIN event_registrations er ON e.id = er.event_id
		WHERE e.id = $1
		GROUP BY e.id, o.name, u.first_name, u.last_name`, eventID).Scan(
		&event.ID, &event.Title, &event.Description, &event.StartDate, &event.EndDate, &event.UserID, &event.Status, &event.Slug, &event.Thumbnail, &event.CreatedAt, &event.UpdatedAt, &event.OrganizationID, &event.MaxRegistration, &event.Organization, &event.Author, &event.TotalRegistered)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// GetEventBySlug Get event and join the total number of registered users
func GetEventBySlug(slug string) (*models.Event, error) {
	var event models.Event
	err := database.DB.QueryRow(context.Background(), `
		SELECT e.id, e.title, e.description, e.start_date, e.end_date, e.user_id, e.status, e.slug, e.thumbnail, e.created_at, e.updated_at, e.organization_id, e.max_registration, o.name as organization, CONCAT(u.first_name, ' ', u.last_name) AS author, COUNT(er.user_id) as total_registered
		FROM events e
		LEFT JOIN organizations o ON e.organization_id = o.id
		LEFT JOIN users u ON e.user_id = u.id
		LEFT JOIN event_registrations er ON e.id = er.event_id
		WHERE e.slug = $1
		GROUP BY e.id, o.name, u.first_name, u.last_name`, slug).Scan(
		&event.ID, &event.Title, &event.Description, &event.StartDate, &event.EndDate, &event.UserID, &event.Status, &event.Slug, &event.Thumbnail, &event.CreatedAt, &event.UpdatedAt, &event.OrganizationID, &event.MaxRegistration, &event.Organization, &event.Author, &event.TotalRegistered)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// ListEvents returns a list of events based on the query parameters
func ListEvents(queryParams map[string]string) ([]*models.Event, int, error) {
	// Limit return data to 10 records per page
	limit := 10

	// Build the query
	query := `
		SELECT e.id, e.title, e.description, e.start_date, e.end_date, e.user_id, e.status, e.slug, e.thumbnail, e.created_at, e.updated_at, e.organization_id, e.max_registration, o.name AS organization, CONCAT(u.first_name, ' ', u.last_name) AS author
		FROM events e
		LEFT JOIN organizations o ON e.organization_id = o.id
		LEFT JOIN users u ON e.user_id = u.id
		WHERE 1 = 1`

	// Add query parameters to the query
	if queryParams["organization_id"] != "" {
		query += " AND o.id = " + queryParams["organization_id"]
	}

	if queryParams["status"] != "" {
		query += " AND e.status = '" + queryParams["status"] + "'"
	}

	var totalRecords int
	err := database.DB.QueryRow(context.Background(), "SELECT COUNT(*) FROM events").Scan(&totalRecords)
	if err != nil {
		return nil, 0, err
	}

	totalPages := (totalRecords + limit - 1) / limit

	if queryParams["page"] != "" {
		page, err := strconv.Atoi(queryParams["page"])
		if err != nil {
			return nil, totalPages, err
		}
		offset := (page - 1) * limit
		query += fmt.Sprintf(" ORDER BY e.created_at DESC LIMIT %d OFFSET %d", limit, offset)
	}

	// Execute the query
	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		return nil, totalPages, err
	}

	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID, &event.Title, &event.Description, &event.StartDate, &event.EndDate, &event.UserID, &event.Status, &event.Slug, &event.Thumbnail, &event.CreatedAt, &event.UpdatedAt, &event.OrganizationID, &event.MaxRegistration, &event.Organization, &event.Author)
		if err != nil {
			return nil, totalPages, err
		}
		events = append(events, &event)
	}

	return events, totalPages, nil
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

	// Check if the user is already registered for the event
	var exists bool
	err = tx.QueryRow(context.Background(), `
			SELECT EXISTS (SELECT 1 FROM event_registrations WHERE event_id = $1 AND user_id = $2)`, eventID, userID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return utils.AlreadyRegisteredError{EventID: eventID}
	}

	_, err = database.DB.Exec(context.Background(), `
        INSERT INTO event_registrations (event_id, user_id, registration_date)
        VALUES ($1, $2, $3)`, eventID, userID, time.Now())
	return err
}

func ListRegisteredUsers(eventID int) ([]*models.User, error) {
	rows, err := database.DB.Query(context.Background(), `
		SELECT u.id, u.username, u.first_name, u.last_name, u.email, u.student_id, u.major, u.profile_picture, u.date_of_birth, u.role_id, u.created_at, u.updated_at, u.year, u.institution_name
		FROM users u
		JOIN event_registrations er ON u.id = er.user_id
		WHERE er.event_id = $1`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.StudentID, &user.Major, &user.ProfilePicture, &user.DateOfBirth, &user.RoleID, &user.CreatedAt, &user.UpdatedAt, &user.Year, &user.InstitutionName)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func ListEventsRegisteredByUser(userID uuid.UUID) ([]*models.Event, error) {
	rows, err := database.DB.Query(context.Background(), `
		SELECT e.id, e.title, e.description, e.start_date, e.end_date, e.user_id, e.status, e.slug, e.thumbnail, e.created_at, e.updated_at, e.organization_id, e.max_registration, o.name as organization_name
		FROM events e
		JOIN event_registrations er ON e.id = er.event_id
		JOIN organizations o ON e.organization_id = o.id
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
			&event.ID, &event.Title, &event.Description, &event.StartDate, &event.EndDate, &event.UserID, &event.Status, &event.Slug, &event.Thumbnail, &event.CreatedAt, &event.UpdatedAt, &event.OrganizationID, &event.MaxRegistration, &event.Organization)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}

func TotalRegisteredUsers(eventID int) (int, error) {
	var totalRegistered int
	err := database.DB.QueryRow(context.Background(), `
		SELECT COUNT(*) FROM event_registrations WHERE event_id = $1`, eventID).Scan(&totalRegistered)
	if err != nil {
		return 0, err
	}
	return totalRegistered, nil
}
