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
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "permissions:list")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": []string{"Permission Denied"}})
		return
	}

	permissions, err := h.PermissionService.ListPermission()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Permissions Fetched Successfully",
		"data":    permissions,
	})
}

func (h *Handler) AssignPermissionToRole(c *gin.Context) {
	roleIDStr := c.Param("roleID")
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Invalid Role ID"}})
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "permissions:assign")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": []string{"Permission Denied"}})
		return
	}

	var permissionIDs []int

	if err := c.ShouldBindJSON(&permissionIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	if err := h.PermissionService.AssignPermissionToRole(roleID, permissionIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Permissions Assigned Successfully",
		"data": gin.H{
			"role_id":     roleID,
			"permissions": permissionIDs,
		},
	})
}
