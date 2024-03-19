package event

import (
	"Backend/internal/handlers/auth"
	"Backend/internal/models"
	"Backend/internal/services"
	"Backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
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
	userID, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "events:create")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	var newEvent models.Event
	if err := c.BindJSON(&newEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	newEvent.UserID = userID

	if newEvent.Title != "" {
		newEvent.Link = "/events/" + utils.GenerateFriendlyURL(newEvent.Title)
	}

	// Check if start date is before end date
	if newEvent.StartDate.After(newEvent.EndDate) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Start Date cannot be after End Date"}})
		return
	}

	if err := h.EventService.CreateEvent(&newEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Event Created Successfully",
		"data":    newEvent,
		"relationships": gin.H{
			"author": gin.H{
				"id": userID,
			},
		},
	})
}

func (h *Handlers) EditEvent(c *gin.Context) {
	userID, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "events:edit")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	eventIDStr := c.Param("eventID")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Invalid Event ID"}})
		return
	}

	var updatedEvent models.Event
	if err := c.BindJSON(&updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	if updatedEvent.StartDate.After(updatedEvent.EndDate) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Start Date cannot be after End Date"}})
		return
	}

	existingEvent, err := h.EventService.GetEventByID(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	if updatedEvent.Title != "" && updatedEvent.Title != existingEvent.Title {
		updatedEvent.Link = "/events/" + utils.GenerateFriendlyURL(updatedEvent.Title)
	} else {
		updatedEvent.Link = existingEvent.Link
	}

	utils.ReflectiveUpdate(existingEvent, updatedEvent)

	if err := h.EventService.EditEvent(eventID, existingEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Event Updated Successfully",
		"data":    existingEvent,
		"relationships": gin.H{
			"author": gin.H{
				"id": userID,
			},
		},
	})
}

func (h *Handlers) DeleteEvent(c *gin.Context) {
	_, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "events:delete")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	eventIDStr := c.Param("eventID")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Invalid Event ID"}})
		return
	}

	if err := h.EventService.DeleteEvent(eventID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Event Deleted Successfully",
	})
}

func (h *Handlers) GetEventByID(c *gin.Context) {
	eventIDStr := c.Param("eventID")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Invalid Event ID"}})
		return
	}

	event, err := h.EventService.GetEventByID(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Event Retrieved Successfully",
		"data":    event,
	})
}

func (h *Handlers) ListEvents(c *gin.Context) {
	log.Println("List Events Begin")

	queryParams := map[string]string{
		"organization_id": c.Query("organization_id"),
		"status":          c.Query("status"),
	}

	events, err := h.EventService.ListEvents(queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	log.Println("List Events End")

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"totalResults": len(events),
		"data":         events,
	})
}

func (h *Handlers) RegisterForEvent(c *gin.Context) {
	userID, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "events:register")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	eventIDStr := c.Param("eventID")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Invalid Event ID"}})
		return
	}

	if err := h.EventService.RegisterForEvent(userID, eventID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Registered Successfully",
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
	userID, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "events:listRegisteredUsers")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	eventIDStr := c.Param("eventID")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Invalid Event ID"}})
		return
	}

	users, err := h.EventService.ListRegisteredUsers(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Registered Users Retrieved Successfully",
		"data":    users,
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

func (h *Handlers) ListEventsRegisteredByUser(c *gin.Context) {
	userID, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "users:edit")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	events, err := h.EventService.ListEventsRegisteredByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Registered Events Retrieved Successfully",
		"data":    events,
	})
}
