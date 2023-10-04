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

	if err := h.UserService.RegisterUser(&newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	log.Println("After calling CreateUser")

	c.JSON(http.StatusCreated, gin.H{
		"data": newUser,
		"meta": gin.H{"message": "User Registered Successfully"},
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
		"data": user,
		"meta": gin.H{"token": token},
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
	c.JSON(http.StatusOK, gin.H{"meta": gin.H{"message": "Logged out successfully"}})
}

func (h *Handlers) GetUserByID(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "view:users")
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

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *Handlers) EditUser(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "edit:users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"Permission Denied"}})
		return
	}

	var updatedUser models.User
	if err := c.BindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	if err := h.UserService.EditUser(userID, &updatedUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedUser})
}

func (h *Handlers) DeleteUser(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	if err := h.UserService.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "User deleted successfully"}})
}

func (h *Handlers) ListUsers(c *gin.Context) {
	users, err := h.UserService.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
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
		"data": gin.H{"message": "Token refreshed successfully"},
		"meta": gin.H{"token": token},
	})
}
