package permission

import (
	"Backend/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handler struct {
	PermissionService *services.PermissionService
}

func NewPermissionHandler(permissionService *services.PermissionService) *Handler {
	return &Handler{PermissionService: permissionService}
}

func (h *Handler) ListPermissions(c *gin.Context) {
	permissions, err := h.PermissionService.ListPermission()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": permissions})
}

func (h *Handler) AssingPermissionToRole(c *gin.Context) {
	roleIDStr := c.Param("roleID")
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Invalid Role ID"}})
		return
	}

	var permissionIDs []int
	if err := c.BindJSON(&permissionIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	if err := h.PermissionService.AssignPermissionToRole(roleID, permissionIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"meta": gin.H{"message": "Permission assigned to role successfully"}})
}
