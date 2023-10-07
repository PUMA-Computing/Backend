package event

import (
	"Backend/internal/models"
	"Backend/internal/services"
	"Backend/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Handlers struct {
	EventService      *services.EventService
	PermissionService *services.PermissionService
}

func NewEventHandlers(eventService *services.EventService, permissionService *services.PermissionService) *Handlers {
	return &Handlers{
		EventService:      eventService,
		PermissionService: permissionService,
	}
}

func (h *Handlers) CreateEvent(c *gin.Context) {
	token, err := utils.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	userID, err := utils.GetUserIDFromToken(token, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "events:create")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"Permission Denied"}})
		return
	}

	var newEvent models.Event
	if err := c.BindJSON(&newEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	now := time.Now()
	if newEvent.StartDate.Before(now) {
		newEvent.Status = "Upcoming"
	} else if newEvent.StartDate.After(now) && newEvent.EndDate.Before(now) {
		newEvent.Status = "Open"
	} else {
		newEvent.Status = "Completed"
	}

	newEvent.Link = "/events/" + utils.GenerateFriendlyURL(newEvent.Title)

	if newEvent.Thumbnail == "" {
		newEvent.Thumbnail = "default.jpg"
	}

	newEvent.UserID = userID
	newEvent.CreatedAt = time.Time{}
	newEvent.UpdatedAt = time.Time{}

	if err := h.EventService.CreateEvent(&newEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"types":      "events",
			"attributes": newEvent,
		},
		"relationships": gin.H{
			"author": gin.H{
				"data": gin.H{
					"id": userID,
				},
			},
		},
	})
}

func (h *Handlers) EditEvent(c *gin.Context) {
	token, err := utils.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	userID, err := utils.GetUserIDFromToken(token, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "events:edit")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"You don't have permission to edit events"}})
		return
	}

	eventIDStr := c.Param("eventID")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Invalid Event ID"}})
		return
	}

	// Check if the user is the author of the event

	//if userID != existingEvent.UserID {
	//	c.JSON(http.StatusForbidden, gin.H{"errors": []string{"You don't have permission to edit this event"}})
	//	return
	//}

	var updatedEvent models.Event
	if err := c.BindJSON(&updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	updatedEvent.UpdatedAt = time.Time{}

	updatedAttributes := make(map[string]interface{})

	if updatedEvent.Title != "" {
		updatedAttributes["title"] = updatedEvent.Title
	}

	if updatedEvent.Description != "" {
		updatedAttributes["description"] = updatedEvent.Description
	}

	if updatedEvent.StartDate.IsZero() {
		updatedAttributes["start_date"] = updatedEvent.StartDate
	}

	if updatedEvent.EndDate.IsZero() {
		updatedAttributes["end_date"] = updatedEvent.EndDate
	}

	if updatedEvent.Status != "" {
		updatedAttributes["status"] = updatedEvent.Status
	}

	if updatedEvent.Thumbnail != "" {
		updatedAttributes["thumbnail"] = updatedEvent.Thumbnail
	}
	if err := h.EventService.EditEvent(eventID, &updatedEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"type":       "events",
			"id":         eventID,
			"attributes": updatedAttributes,
		},
		"relationships": gin.H{
			"author": gin.H{
				"data": gin.H{
					"id": userID,
				},
			},
		},
	})
}

func (h *Handlers) DeleteEvent(c *gin.Context) {
	token, err := utils.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	userID, err := utils.GetUserIDFromToken(token, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "events:delete")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"Permission Denied"}})
		return
	}

	eventIDStr := c.Param("eventID")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Invalid Event ID"}})
		return
	}

	if err := h.EventService.DeleteEvent(eventID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "Event Deleted Successfully"}})
}

func (h *Handlers) GetEventByID(c *gin.Context) {
	eventIDStr := c.Param("eventID")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Invalid Event ID"}})
		return
	}

	event, err := h.EventService.GetEventByID(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"type":       "events",
			"id":         event.ID,
			"attributes": event,
		},
		"relationships": gin.H{
			"author": gin.H{
				"data": gin.H{
					"id": event.UserID,
				},
			},
		},
	})
}

func (h *Handlers) ListEvents(c *gin.Context) {
	log.Println("List Events Begin")
	events, err := h.EventService.ListEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	log.Println("List Events End")

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"type":       "events",
			"attributes": events,
		},
	})
}

func (h *Handlers) RegisterForEvent(c *gin.Context) {
	token, err := utils.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	userID, err := utils.GetUserIDFromToken(token, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "events:register")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"Permission Denied"}})
		return
	}

	eventIDStr := c.Param("eventID")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Invalid Event ID"}})
		return
	}

	if err := h.EventService.RegisterForEvent(userID, eventID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"type": "events_registration",
			"id":   eventID,
			"attributes": gin.H{
				"user_id":  userID,
				"event_id": eventID,
			},
		},
		"relationships": gin.H{
			"user": gin.H{
				"data": gin.H{
					"id": userID,
				},
			},
			"event": gin.H{
				"data": gin.H{
					"id": eventID,
				},
			},
		},
	})
}

func (h *Handlers) ListRegisteredUsers(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "events:listRegisteredUsers")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"Permission Denied"}})
		return
	}

	eventIDStr := c.Param("eventID")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Invalid Event ID"}})
		return
	}

	users, err := h.EventService.ListRegisteredUsers(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"type":       "events_registration",
			"id":         eventID,
			"attributes": users,
		},
		"relationships": gin.H{
			"user": gin.H{
				"data": gin.H{
					"id": userID,
				},
			},
		},
	})
}
