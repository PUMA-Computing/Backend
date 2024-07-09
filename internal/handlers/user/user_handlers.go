package user

import (
	"Backend/internal/handlers/auth"
	"Backend/internal/models"
	"Backend/internal/services"
	"Backend/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Handlers struct {
	UserService       *services.UserService
	PermissionService *services.PermissionService
	AWSService        *services.S3Service
	R2Service         *services.S3Service
}

func NewUserHandlers(userService *services.UserService, permissionService *services.PermissionService, awsService *services.S3Service, r2Service *services.S3Service) *Handlers {
	return &Handlers{
		UserService:       userService,
		PermissionService: permissionService,
		AWSService:        awsService,
		R2Service:         r2Service,
	}
}

func (h *Handlers) GetUserByID(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Access token is expired or invalid"})
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

// EditUser User can only edit their own profile
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

	// Check if password is empty
	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}

		updatedAttributes["password"] = string(hashedPassword)
		updatedUser.Password = string(hashedPassword)
	}

	if updatedUser.FirstName != "" {
		updatedAttributes["first_name"] = updatedUser.FirstName
	}

	if updatedUser.MiddleName != nil {
		updatedAttributes["middle_name"] = updatedUser.MiddleName
	}

	if updatedUser.LastName != "" {
		updatedAttributes["last_name"] = updatedUser.LastName
	}

	if updatedUser.Email != "" {
		updatedAttributes["email"] = updatedUser.Email
	}

	if updatedUser.StudentID != "" {
		// Check if student ID already exists
		studentIDExists, err := h.UserService.CheckStudentIDExists(updatedUser.StudentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}

		if studentIDExists {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Student ID already exists"}})
			return
		}

		// Check if student ID is valid
		if len(updatedUser.StudentID) != 12 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Student ID must be 12 characters long"}})
			return
		}
		if updatedUser.StudentID[:3] != "001" && updatedUser.StudentID[:3] != "012" && updatedUser.StudentID[:3] != "013" && updatedUser.StudentID[:3] != "025" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"You are not a student of faculty of computing"}})
			return
		}

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

	if updatedUser.ProfilePicture != "" {
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
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{"Unauthorized"}})
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

func (h *Handlers) GetRoleIDByUserID(c *gin.Context, userID uuid.UUID) (int, error) {
	roleID, err := h.UserService.GetRoleIDByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return 0, err
	}

	return roleID, nil
}

func (h *Handlers) AdminUpdateRoleAndStudentIDVerified(c *gin.Context) {
	_, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "users:create")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	userIDStr := c.Param("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Invalid User ID"}})
		return
	}

	// Bind JSON
	var updatedUser models.User
	if err := c.BindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	if err := h.UserService.AdminUpdateRoleAndStudentIDVerified(userID, updatedUser.RoleID, updatedUser.StudentIDVerified); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User Updated Successfully",
	})
}

func (h *Handlers) UploadProfilePicture(c *gin.Context) {
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

	file, _, err := c.Request.FormFile("profile_picture")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	// Check if file size is greater than 2MB
	optimizedImage, err := utils.OptimizeImage(file, 2800, 1080)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	optimizedImageBytes, err := io.ReadAll(optimizedImage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	// Choose storage service to upload image to (AWS or R2)
	upload := utils.ChooseStorageService()

	// Upload image to storage service
	if upload == utils.R2Service {
		err = h.R2Service.UploadFileToR2(context.Background(), "users", user.ID.String(), optimizedImageBytes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}

		user.ProfilePicture, _ = h.R2Service.GetFileR2("users", userID.String())
	} else {
		err = h.AWSService.UploadFileToAWS(context.Background(), "users", user.ID.String(), optimizedImageBytes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}

		user.ProfilePicture, _ = h.AWSService.GetFileAWS("users", user.ID.String())
	}

	// Update user profile picture
	if err := h.UserService.UploadProfilePicture(userID, user.ProfilePicture); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile Picture Uploaded Successfully",
		"data":    user.ProfilePicture,
	})
}

func (h *Handlers) EnableTwoFA(c *gin.Context) {
	userID, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "users:2fa")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	qr, setupKey, err := h.UserService.EnableTwoFA(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Two Factor Authentication Enabled Successfully",
		"data": gin.H{
			"qr_image":  qr,
			"setup_key": setupKey,
		},
	})
}

func (h *Handlers) VerifyTwoFA(c *gin.Context) {
	userID, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "users:2fa")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	var request struct {
		Code string `json:"code"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	valid, err := h.UserService.VerifyTwoFA(userID, request.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid TOTP Code"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "TOTP Verified Successfully"})
}

func (h *Handlers) DisableTwoFA(c *gin.Context) {
	userID, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "users:2fa")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	err = h.UserService.DisableTwoFA(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Two Factor Authentication Disabled Successfully"})
}
