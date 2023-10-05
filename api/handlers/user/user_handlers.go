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

func (h *Handlers) RegisterUser(c *gin.Context) {
	var newUser models.User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	log.Println("Before calling CreateUser")

	newUser.RoleID = 2
	newUser.CreatedAt = time.Time{}
	newUser.UpdatedAt = time.Time{}

	err := h.UserService.RegisterUser(&newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	log.Println("After calling CreateUser")

	c.JSON(http.StatusCreated, gin.H{
		"jsonapi": gin.H{
			"version": "1.1",
		},
		"data": gin.H{
			"type":       "users",
			"attributes": newUser,
		},
		"relationships": gin.H{
			"role": gin.H{
				"data": gin.H{
					"type": "roles",
					"id":   newUser.RoleID,
				},
			},
		},
	})
}

func (h *Handlers) Login(c *gin.Context) {
	var loginRequest models.User
	if err := c.BindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	user, err := h.UserService.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	token, err := utils.GenerateJWTToken(user.ID, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	if err := utils.StoreTokenInRedis(user.ID, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jsonapi": gin.H{
			"version": "1.1",
		},
		"data": gin.H{
			"id":   user.ID,
			"type": "token",
			"attributes": gin.H{
				"access_token": token,
				"token_type":   "Bearer",
				"expires_in":   86400,
			},
		},
	})
}

func (h *Handlers) Logout(c *gin.Context) {
	tokenString, err := utils.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": []string{"Unauthorized"}})
		return
	}

	_, err = utils.ValidateToken(tokenString, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": []string{"Unauthorized"}})
		return
	}

	err = utils.RevokeToken(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "Logout Successful"}})
}

func (h *Handlers) GetUserByID(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "get:users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"Permission Denied"}})
		return
	}

	user, err := h.UserService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": []string{"User not found"}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"type":       "users",
			"id":         userID,
			"attributes": user,
		},
		"relationships": gin.H{
			"role": gin.H{
				"data": gin.H{
					"type": "roles",
					"id":   user.RoleID,
				},
			},
		},
	})
}

func (h *Handlers) EditUser(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	log.Println("userID: ", userID)

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "users:edit")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	log.Println("hasPermission: ", hasPermission)

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"Permission Denied"}})
		return
	}

	log.Println("Before binding JSON")

	var updatedUser models.User
	if err := c.BindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
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

	log.Println("After binding JSON")

	if err := h.UserService.EditUser(userID, &updatedUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	log.Println("After calling EditUser")

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"type":       "users",
			"id":         userID,
			"attributes": updatedAttributes,
		},
	})
}

func (h *Handlers) DeleteUser(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "users:delete")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"Permission Denied"}})
		return
	}

	if err := h.UserService.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "User Deleted Successfully"}})
}

func (h *Handlers) ListUsers(c *gin.Context) {
	log.Println("Before calling GetUserIDFromContext")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	log.Println("userID: ", userID)

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "users:list")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	log.Println("hasPermission: ", hasPermission)

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"Permission Denied"}})
		return
	}

	log.Println("Before calling ListUsers")

	users, err := h.UserService.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	log.Println("After calling ListUsers")

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"type":       "users",
			"attributes": users,
		}})
}

func (h *Handlers) RefreshToken(c *gin.Context) {
	tokenString, err := utils.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": []string{"Unauthorized"}})
		return
	}

	claims, err := utils.ValidateToken(tokenString, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": []string{"Unauthorized"}})
		return
	}

	userID := claims.UserID
	token, err := utils.GenerateJWTToken(userID, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	if err := utils.StoreTokenInRedis(userID, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"types": "Bearer",
			"attributes": gin.H{
				"access_token": token,
			},
		},
	})
}
