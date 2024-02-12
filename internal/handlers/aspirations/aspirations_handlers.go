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

	if err := h.AspirationService.CreateAspiration(&newAspiration); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": []string{"Aspiration created successfully"}})
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
