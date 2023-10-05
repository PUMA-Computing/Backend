package permission

import (
	"Backend/internal/services"
	"Backend/pkg/utils"
	"context"
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
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "permissions:list")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"Permission Denied"}})
		return
	}

	permissions, err := h.PermissionService.ListPermission()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": permissions})
}

func (h *Handler) AssignPermissionToRole(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "permissions:assign")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"Permission Denied"}})
		return
	}

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
