package aspirations

import (
	"Backend/internal/models"
	"Backend/internal/services"
	"Backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
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
	userID, err := utils.ExtractUserIDAndCheckPermission(c, "aspirations:create")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	var aspiration models.Aspiration
	if err := c.ShouldBindJSON(&aspiration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	aspiration.UserID = userID
	if err := h.AspirationService.CreateAspiration(&aspiration); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": []string{"Aspiration created successfully"}})
}
