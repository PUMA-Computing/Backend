package user

import (
	"Backend/internal/models"
	"Backend/internal/services"
	"Backend/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

type Handlers struct {
	UserService       *services.UserService
	PermissionService *services.PermissionService
}

func NewUserHandlers(userService *services.UserService, permissionService *services.PermissionService) *Handlers {
	return &Handlers{
		UserService:       userService,
		PermissionService: permissionService,
	}
}

func (h *Handlers) GetUserByID(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	user, err := h.UserService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User Retrieved Successfully",
		"data":    user,
	})
}

func (h *Handlers) EditUser(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	log.Println("userID: ", userID)

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "users:edit")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	log.Println("hasPermission: ", hasPermission)

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": []string{"Permission Denied"}})
		return
	}

	log.Println("Before binding JSON")

	var updatedUser models.User
	if err := c.BindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	updatedUser.RoleID = 2
	updatedUser.UpdatedAt = time.Time{}

	updatedAttributes := make(map[string]interface{})

	if updatedUser.Username != "" {
		updatedAttributes["username"] = updatedUser.Username
	}

	if updatedUser.Password != "" {
		updatedAttributes["password"] = updatedUser.Password
	}

	if updatedUser.FirstName != "" {
		updatedAttributes["first_name"] = updatedUser.FirstName
	}

	if updatedUser.MiddleName != "" {
		updatedAttributes["middle_name"] = updatedUser.MiddleName
	}

	if updatedUser.LastName != "" {
		updatedAttributes["last_name"] = updatedUser.LastName
	}

	if updatedUser.Email != "" {
		updatedAttributes["email"] = updatedUser.Email
	}

	if updatedUser.StudentID != "" {
		updatedAttributes["student_id"] = updatedUser.StudentID
	}

	if updatedUser.Major != "" {
		updatedAttributes["major"] = updatedUser.Major
	}

	if updatedUser.Year != "" {
		updatedAttributes["year"] = updatedUser.Year
	}

	if updatedUser.DateOfBirth != nil {
		updatedAttributes["date_of_birth"] = updatedUser.DateOfBirth
	}

	if updatedUser.ProfilePicture != nil {
		updatedAttributes["profile_picture"] = updatedUser.ProfilePicture
	}

	log.Println("After binding JSON")

	if err := h.UserService.EditUser(userID, &updatedUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	log.Println("After calling EditUser")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User Updated Successfully",
		"data":    updatedAttributes,
	})
}

func (h *Handlers) DeleteUser(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "users:delete")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": []string{"Permission Denied"}})
		return
	}

	if err := h.UserService.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User Deleted Successfully",
	})
}

func (h *Handlers) ListUsers(c *gin.Context) {
	tokenString, err := utils.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": []string{"Unauthorized"}})
		return
	}

	claims, err := utils.ValidateToken(tokenString, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": []string{"Unauthorized"}})
		return
	}

	userID := claims.UserID

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "users:list")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": []string{"Permission Denied"}})
		return
	}

	log.Println("Before calling ListUsers")
	users, err := h.UserService.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Users Retrieved Successfully",
		"data":    users,
	})
}
