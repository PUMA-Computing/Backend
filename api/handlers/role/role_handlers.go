package role

import (
	"Backend/internal/models"
	"Backend/internal/services"
	"Backend/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type Handler struct {
	roleService       *services.RoleService
	UserService       *services.UserService
	PermissionService *services.PermissionService
}

func NewRoleHandler(roleService *services.RoleService, userService *services.UserService, permissionService *services.PermissionService) *Handler {
	return &Handler{
		roleService:       roleService,
		UserService:       userService,
		PermissionService: permissionService,
	}
}

func (h *Handler) CreateRole(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "roles:create")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"Permission Denied"}})
		return
	}

	var newRole models.Roles
	if err := c.BindJSON(&newRole); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	if err := h.roleService.CreateRole(&newRole); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": newRole,
		"meta": gin.H{"message": "Role Created Successfully"},
	})
}

func (h *Handler) GetRoleByID(c *gin.Context) {
	roleIDStr := c.Param("roleID")
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Invalid Role ID"}})
		return
	}

	role, err := h.roleService.GetRoleByID(roleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": []string{"Role not found"}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": role})
}

func (h *Handler) EditRole(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "roles:edit")
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

	var updatedRole models.Roles
	if err := c.BindJSON(&updatedRole); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	if err := h.roleService.EditRole(roleID, &updatedRole); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": updatedRole,
		"meta": gin.H{"message": "Role Updated Successfully"},
	})
}

func (h *Handler) DeleteRole(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "roles:delete")
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

	if err := h.roleService.DeleteRole(roleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"meta": gin.H{"message": "Role Deleted Successfully"}})
}

func (h *Handler) ListRoles(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "roles:list")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"Permission Denied"}})
		return
	}

	roles, err := h.roleService.ListRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": roles})
}

func (h *Handler) AssignRoleToUser(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "roles:assign")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"Permission Denied"}})
		return
	}

	targetUserIDStr := c.Param("userID")
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Invalid User ID"}})
		return
	}

	targetUser, err := h.UserService.GetUserByID(targetUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": []string{"User not found"}})
		return
	}

	if targetUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": []string{"User not found"}})
		return
	}

	roleIDStr := c.Param("roleID")
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Invalid Role ID"}})
		return
	}

	if err := h.roleService.AssignRoleToUser(targetUserID, roleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"meta": gin.H{"message": "Role Assigned Successfully"}})
}
