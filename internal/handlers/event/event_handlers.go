package event

import (
	"Backend/internal/handlers/auth"
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
	S3Service         *services.S3Service
}

func NewEventHandlers(eventService *services.EventService, permissionService *services.PermissionService, s3Service *services.S3Service) *Handlers {
	return &Handlers{
		EventService:      eventService,
		PermissionService: permissionService,
		S3Service:         s3Service,
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
		newEvent.Slug = "/event/" + utils.GenerateFriendlyURL(newEvent.Title)
	}

	// Check if start date is before end date
	if newEvent.StartDate.After(newEvent.EndDate) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Start Date cannot be after End Date"}})
		return
	}

	// Upload thumbnail to AWS S3 and save the URL
	if newEvent.Thumbnail != "" {
		if err := h.S3Service.UploadFile(context.TODO(), newEvent.Slug+".jpg", []byte(newEvent.Thumbnail)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{"Failed to upload thumbnail to AWS S3"}})
			return
		}
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
	_, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "events:edit")
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

	file, err := c.FormFile("thumbnail")
	if err == nil {
		// If a thumbnail file is provided, upload it to AWS S3
		imageReader, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Failed to open thumbnail file"}})
			return
		}
		defer imageReader.Close()

		fileBytes, err := io.ReadAll(imageReader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{"Failed to read thumbnail file"}})
			return
		}

		fileKey := "event/" + strconv.Itoa(eventID) + ".jpg"
		if err := h.S3Service.UploadFile(context.TODO(), fileKey, fileBytes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}

		// Get the URL of the uploaded file
		fileURL, _ := h.S3Service.GetFile(context.TODO(), fileKey)
		log.Println("File URL: ", fileURL)

		// Read the JSON payload into a map
		var payload map[string]interface{}
		if err := c.BindJSON(&payload); err != nil {
			log.Println("Error binding JSON: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}

		// Replace the thumbnail key in the map
		payload["thumbnail"] = fileURL

		// Convert the map back to JSON
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			log.Println("Error marshalling JSON: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}

		// Read the JSON into the updatedEvent object
		var updatedEvent models.Event
		if err := json.Unmarshal(payloadBytes, &updatedEvent); err != nil {
			log.Println("Error unmarshalling JSON: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}

		// Check if start date is before end date
		if updatedEvent.StartDate.After(updatedEvent.EndDate) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Start Date cannot be after End Date"}})
			return
		}

		if err := h.EventService.EditEvent(eventID, &updatedEvent); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Event Updated Successfully",
			"data":    updatedEvent,
		})
	} else {
		// If no thumbnail file is provided, update the event without changing the thumbnail
		var updatedEvent models.Event
		if err := c.BindJSON(&updatedEvent); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}

		// Check if start date is before end date
		if updatedEvent.StartDate.After(updatedEvent.EndDate) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Start Date cannot be after End Date"}})
			return
		}

		if err := h.EventService.EditEvent(eventID, &updatedEvent); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Event Updated Successfully",
			"data":    updatedEvent,
		})
	}

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
		"slug":            c.Query("slug"),
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
	log.Println("Register for Event Begin")
	userID, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "events:register")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
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
