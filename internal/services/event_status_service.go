package services

import (
	"Backend/internal/models"
	"log"
	"time"
)

type EventStatusUpdater struct {
	EventService *EventService
}

func NewEventStatusUpdater(eventService *EventService) *EventStatusUpdater {
	return &EventStatusUpdater{EventService: eventService}
}

func (e *EventStatusUpdater) Run() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("EventStatusUpdater: Status update loop started")
			// Get all events
			events, err := e.EventService.ListEvents(map[string]string{})
			if err != nil {
				log.Printf("EventStatusUpdater: Error fetching events: %v\n", err)
				continue
			}

			// Update event status
			for _, event := range events {
				// Update status based on current time
				e.updateEventStatus(event)
			}
		}
	}
}

func (e *EventStatusUpdater) updateEventStatus(event *models.Event) {
	currentTime := time.Now()

	switch {
	case currentTime.Before(event.StartDate):
		event.Status = "Upcoming"
	case currentTime.After(event.StartDate) && currentTime.Before(event.EndDate):
		event.Status = "Open"
	default:
		event.Status = "Ended"
	}

	// Update the event status in the database
	if err := e.EventService.EditEvent(event.ID, event); err != nil {
		// Log the error
		log.Println("Error updating event status:", err)
	}
}
