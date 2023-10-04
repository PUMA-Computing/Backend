package role

import (
	"Backend/internal/models"
	"Backend/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handler struct {
	roleService *services.RoleService
}

func NewRoleHandler(roleService *services.RoleService) *Handler {
	return &Handler{
		roleService: roleService,
	}
}

func (h *Handler) CreateRole(c *gin.Context) {
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
	roles, err := h.roleService.ListRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": roles})
}

func (h *Handler) AssignRoleToUser(c *gin.Context) {
	roleIDStr := c.Param("roleID")
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Invalid Role ID"}})
		return
	}

	userIDStr := c.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Invalid User ID"}})
		return
	}

	if err := h.roleService.AssignRoleToUser(userID, roleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"meta": gin.H{"message": "Role Assigned Successfully"}})
}
