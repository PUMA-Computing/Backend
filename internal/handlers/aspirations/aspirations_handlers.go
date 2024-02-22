package aspirations

import (
	"Backend/internal/handlers/auth"
	"Backend/internal/models"
	"Backend/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handlers struct {
	AspirationService *services.AspirationService
	PermissionService *services.PermissionService
}

func NewAspirationHandlers(aspirationService *services.AspirationService, permissionService *services.PermissionService) *Handlers {
	return &Handlers{
		AspirationService: aspirationService,
		PermissionService: permissionService,
	}
}

func (h *Handlers) CreateAspiration(c *gin.Context) {
	userID, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "aspirations:create")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	var newAspiration models.Aspiration
	if err := c.BindJSON(&newAspiration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	newAspiration.UserID = userID

	createdAspiration, err := h.AspirationService.CreateAspiration(&newAspiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "aspiration": createdAspiration})
}

func (h *Handlers) CloseAspiration(c *gin.Context) {
	_, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "aspirations:close")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Invalid ID"}})
		return
	}

	if err := h.AspirationService.CloseAspirationByID(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": []string{"Aspiration closed successfully"}})
}

func (h *Handlers) DeleteAspiration(c *gin.Context) {
	_, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "aspirations:delete")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Invalid ID"}})
		return
	}

	if err := h.AspirationService.DeleteAspirationByID(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": []string{"Aspiration deleted successfully"}})
}

func (h *Handlers) GetAspirations(c *gin.Context) {
	queryParams := map[string]string{
		"organization_id": c.Query("organization_id"),
		"user_id":         c.Query("user_id"),
		"closed":          c.Query("closed"),
		"anonymous":       c.Query("anonymous"),
	}

	aspirations, err := h.AspirationService.GetAspirations(queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	for i, a := range aspirations {
		upvotes, err := h.AspirationService.GetUpvotesByAspirationID(a.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}
		aspirations[i].Upvote = upvotes
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"aspirations":  aspirations,
		"totalResults": len(aspirations),
	})
}

func (h *Handlers) GetAspirationByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Invalid ID"}})
		return
	}

	aspiration, err := h.AspirationService.GetAspirationByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "aspiration": aspiration})
}

func (h *Handlers) UpvoteAspiration(c *gin.Context) {
	userID, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "aspirations:upvote")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Invalid ID"}})
		return
	}

	upvoteExists, err := h.AspirationService.UpvoteExists(userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	if upvoteExists {
		if err := h.AspirationService.RemoveUpvote(userID, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": []string{"Upvote removed successfully"}})
		return
	}

	if err := h.AspirationService.AddUpvote(userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": []string{"Upvote added successfully"}})
}

func (h *Handlers) AddUpvote(c *gin.Context) {
	userID, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "aspirations:upvote")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Invalid ID"}})
		return
	}

	if err := h.AspirationService.AddUpvote(userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": []string{"Upvote added successfully"}})
}

func (h *Handlers) RemoveUpvote(c *gin.Context) {
	userID, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "aspirations:upvote")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Invalid ID"}})
		return
	}

	if err := h.AspirationService.RemoveUpvote(userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": []string{"Upvote removed successfully"}})
}

func (h *Handlers) GetUpvotesByAspirationID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Invalid ID"}})
		return
	}

	upvotes, err := h.AspirationService.GetUpvotesByAspirationID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "upvotes": upvotes})
}
