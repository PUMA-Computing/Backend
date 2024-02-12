package auth

import (
	"Backend/internal/models"
	"Backend/internal/services"
	"Backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

type Handlers struct {
	AuthService       *services.AuthService
	PermissionService *services.PermissionService
}

func NewAuthHandlers(authService *services.AuthService, permissionService *services.PermissionService) *Handlers {
	return &Handlers{
		AuthService:       authService,
		PermissionService: permissionService,
	}
}

func (h *Handlers) RegisterUser(c *gin.Context) {
	var newUser models.User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	if len(newUser.StudentID) != 12 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Student ID must be 12 characters"})
		return
	} else if newUser.StudentID[:3] != "001" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Student ID must start with 001"})
		return
	} else if newUser.StudentID[3:7] < "2010" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Student ID must be no less than 2010"})
		return
	} else if newUser.StudentID[7:] < "0000" || newUser.StudentID[7:] > "9999" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Student ID must be in the format of 001XXXXYYYYY where X is the year and Y is the student number"})
		return
	}

	// Check student id already exists
	_, err := h.AuthService.CheckStudentIDExists(newUser.StudentID)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "Student ID already exists"})
		return
	}

	if newUser.Email[len(newUser.Email)-24:] != "@student.president.ac.id" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Must be President University student email"})
		return
	}

	err = h.AuthService.RegisterUser(&newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User Created Successfully",
	})
}

func (h *Handlers) Login(c *gin.Context) {
	var loginRequest models.User
	if err := c.BindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	user, err := h.AuthService.LoginUser(loginRequest.Username, loginRequest.Password)
	if err != nil {
		log.Println("Error in here")
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	token, err := utils.GenerateJWTToken(user.ID, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	if err := utils.StoreTokenInRedis(user.ID, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login Successful",
		"data": gin.H{
			"user_id": user.ID,
			"type":    "token",
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
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": []string{"Unauthorized"}})
		return
	}

	_, err = utils.ValidateToken(tokenString, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": []string{"Unauthorized"}})
		return
	}

	err = utils.RevokeToken(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Logout Successful"})
}

func (h *Handlers) RefreshToken(c *gin.Context) {
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
	token, err := utils.GenerateJWTToken(userID, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	if err := utils.StoreTokenInRedis(userID, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token Refreshed Successfully",
		"data": gin.H{
			"types": "Bearer",
			"attributes": gin.H{
				"access_token": token,
			},
		},
	})
}
