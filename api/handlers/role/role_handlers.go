package role

import (
	"Backend/internal/models"
	"Backend/internal/services"
	"Backend/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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
	token, err := utils.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	userID, err := utils.GetUserIDFromToken(token, os.Getenv("JWT_SECRET_KEY"))
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

	newRole.CreatedAt = time.Now()
	newRole.UpdatedAt = time.Now()

	if err := h.roleService.CreateRole(&newRole); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"types":      "roles",
			"attributes": newRole,
		},
	})
}

func (h *Handler) GetRoleByID(c *gin.Context) {
	token, err := utils.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	userID, err := utils.GetUserIDFromToken(token, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "roles:get")
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

	role, err := h.roleService.GetRoleByID(roleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": []string{"Role not found"}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"type":       "news",
			"id":         role.ID,
			"attributes": role,
		},
	})
}

func (h *Handler) EditRole(c *gin.Context) {
	token, err := utils.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	userID, err := utils.GetUserIDFromToken(token, os.Getenv("JWT_SECRET_KEY"))
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

	log.Println("Edit Role")

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

	updatedRole.UpdatedAt = time.Now()

	updatedAttributes := make(map[string]interface{})
	if updatedRole.Name != "" {
		updatedAttributes["name"] = updatedRole.Name
	}

	if err := h.roleService.EditRole(roleID, &updatedRole); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"type":       "roles",
			"id":         roleID,
			"attributes": updatedAttributes,
		},
	})
}

func (h *Handler) DeleteRole(c *gin.Context) {
	token, err := utils.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	userID, err := utils.GetUserIDFromToken(token, os.Getenv("JWT_SECRET_KEY"))
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

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"message": "Role Deleted Successfully",
		},
	})
}

func (h *Handler) ListRoles(c *gin.Context) {
	token, err := utils.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	userID, err := utils.GetUserIDFromToken(token, os.Getenv("JWT_SECRET_KEY"))
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

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"type":       "roles",
			"attributes": roles,
		},
	})
}

func (h *Handler) AssignRoleToUser(c *gin.Context) {
	token, err := utils.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	userID, err := utils.GetUserIDFromToken(token, os.Getenv("JWT_SECRET_KEY"))
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

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"type": "users_roles",
			"id":   roleID,
			"attributes": gin.H{
				"user_id": targetUserID,
				"role_id": roleID,
			},
		},
		"relationships": gin.H{
			"author": gin.H{
				"data": gin.H{
					"id": userID,
				},
			},
		},
	})
}
