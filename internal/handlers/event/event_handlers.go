package event

import (
	"Backend/internal/handlers/auth"
	"Backend/internal/handlers/user"
	"Backend/internal/models"
	"Backend/internal/services"
	"Backend/pkg/utils"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Handlers struct {
	EventService      *services.EventService
	PermissionService *services.PermissionService
	AWSService        *services.S3Service
	R2Service         *services.S3Service
}

func NewEventHandlers(eventService *services.EventService, permissionService *services.PermissionService, AWSService *services.S3Service, R2Service *services.S3Service) *Handlers {
	return &Handlers{
		EventService:      eventService,
		PermissionService: permissionService,
		AWSService:        AWSService,
		R2Service:         R2Service,
	}
}

func (h *Handlers) CreateEvent(c *gin.Context) {
	userID, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "events:create")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	data := c.Request.FormValue("data")
	var newEvent models.Event
	if err := json.Unmarshal([]byte(data), &newEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	newEvent.UserID = userID

	if newEvent.Title != "" {
		newEvent.Slug = utils.GenerateFriendlyURL(newEvent.Title)
	}

	if newEvent.StartDate.After(newEvent.EndDate) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Start Date cannot be after End Date"}})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	optimizedImage, err := utils.OptimizeImage(file, 2800, 1080)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	optimizedImageBytes, err := io.ReadAll(optimizedImage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	// Choose storage service to upload image to (AWS or R2)
	upload := utils.ChooseStorageService()

	// Upload image to storage service
	if upload == utils.R2Service {
		err = h.R2Service.UploadFileToR2(context.Background(), "event", newEvent.Slug, optimizedImageBytes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}

		newEvent.Thumbnail, _ = h.R2Service.GetFileR2("event", newEvent.Slug)
	} else {
		err = h.AWSService.UploadFileToAWS(context.Background(), "event", newEvent.Slug, optimizedImageBytes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}

		newEvent.Thumbnail, _ = h.AWSService.GetFileAWS("event", newEvent.Slug)
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

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	data := c.Request.FormValue("data")
	var updatedEvent models.Event
	if err := json.Unmarshal([]byte(data), &updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	if updatedEvent.StartDate.After(updatedEvent.EndDate) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Start Date cannot be after End Date"}})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	optimizedImage, err := utils.OptimizeImage(file, 2800, 1080)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	optimizedImageBytes, err := io.ReadAll(optimizedImage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	// Choose storage service to upload image
	upload := utils.ChooseStorageService()

	// Upload image to storage service
	if upload == utils.R2Service {
		err = h.R2Service.UploadFileToR2(context.Background(), "event", updatedEvent.Slug, optimizedImageBytes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}

		updatedEvent.Thumbnail, _ = h.R2Service.GetFileR2("event", updatedEvent.Slug)
	} else {
		err = h.AWSService.UploadFileToAWS(context.Background(), "event", updatedEvent.Slug, optimizedImageBytes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}

		updatedEvent.Thumbnail, _ = h.AWSService.GetFileAWS("event", updatedEvent.Slug)
	}

	existingEvent, err := h.EventService.GetEventByID(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	if updatedEvent.Title != "" && updatedEvent.Title != existingEvent.Title {
		updatedEvent.Slug = utils.GenerateFriendlyURL(updatedEvent.Title)
	} else {
		updatedEvent.Slug = existingEvent.Slug
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

	event, err := h.EventService.GetEventByID(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	// Check image exists on AWS or R2
	exists, _ := h.AWSService.FileExists(context.Background(), "event", event.Slug)
	if exists {
		if err := h.AWSService.DeleteFile(context.Background(), "event", event.Slug); err != nil {
			log.Println("Error deleting file from AWS")
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}
	} else {
		exists, _ := h.R2Service.FileExists(context.Background(), "event", event.Slug)
		if exists {
			if err := h.R2Service.DeleteFile(context.Background(), "event", event.Slug); err != nil {
				log.Println("Error deleting file from R2")
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
				return
			}
		}
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

// GetEventByID retrieves an event by its ID
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

// GetEventBySlug retrieves an event by its slug
func (h *Handlers) GetEventBySlug(c *gin.Context) {
	slug := c.Param("eventID")

	event, err := h.EventService.GetEventBySlug(slug)
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

// ListEvents retrieves a list of events based on the query parameters
func (h *Handlers) ListEvents(c *gin.Context) {
	log.Println("List Events Begin")

	queryParams := map[string]string{
		"organization_id": c.Query("organization_id"),
		"status":          c.Query("status"),
		"page":            c.Query("page"),
	}

	events, totalPages, err := h.EventService.ListEvents(queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"data":         events,
		"totalResults": len(events),
		"totalPages":   totalPages,
	})
}

func (h *Handlers) RegisterForEvent(c *gin.Context) {
	log.Println("Register for Event Begin")
	userID, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "events:register")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	// Check Role id and if it is 6 cannot register for event
	roleID, err := (&user.Handlers{}).GetRoleIDByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	if roleID == 8 {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": []string{"You are not eligible to this event"}})
		return
	}

	log.Println("Register for Event Middle")

	eventIDStr := c.Param("eventID")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Invalid Event ID"}})
		return
	}

	log.Println("Register for Event Middle 2")

	if err := h.EventService.RegisterForEvent(userID, eventID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	log.Println("Register for Event End")

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
