package auth

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
	log.Println("before register user")
	var newUser models.User
	suffix := "@student.president.ac.id"
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	// Log the user data
	log.Printf("User: %v", newUser)

	if len(newUser.StudentID) == 0 || newUser.StudentID[:3] != "001" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Student ID must start with 001"})
		return
	} else if newUser.StudentID[3:7] < "2010" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Student ID must be no less than 2010"})
		return
	} else if newUser.StudentID[7:] < "0000" || newUser.StudentID[7:] > "9999" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Student ID must be in the format of 001XXXXYYYYY where X is the year and Y is the student number"})
		return
	}

	log.Println("before check student id exists")
	// Check student id already exists
	_, err := h.AuthService.CheckStudentIDExists(newUser.StudentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	log.Println("before check Username or Email exists")
	// Check username already exists
	_, err = h.AuthService.CheckUsernameOrEmailExists(newUser.Username, newUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	// email must be president university student email
	if len(newUser.Email) < len(suffix) || newUser.Email[len(newUser.Email)-len(suffix):] != suffix {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Email must be a President University student email"})
		return
	}

	log.Println("start register user")
	err = h.AuthService.RegisterUser(&newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	log.Println("after register user")

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

	log.Println("Co INi kenapa ajg")
	user, err := h.AuthService.LoginUser(loginRequest.Username, loginRequest.Password)
	if err != nil {
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

func (h *Handlers) ExtractUserIDAndCheckPermission(c *gin.Context, permissionType string) (uuid.UUID, error) {
	token, err := utils.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return uuid.UUID{}, err
	}

	userID, err := utils.GetUserIDFromToken(token, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return uuid.UUID{}, err
	}

	hasPermission, err := (&services.PermissionService{}).CheckPermission(context.Background(), userID, permissionType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return uuid.UUID{}, err
	}

	if !hasPermission {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": []string{"Unauthorized"}})
		return uuid.UUID{}, err
	}

	return userID, nil
}
