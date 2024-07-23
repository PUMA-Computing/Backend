package services

import (
	"Backend/internal/database/app"
	"Backend/internal/models"
	"Backend/pkg/utils"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type UserService struct {
	otp *OTPManager
}

func NewUserService() *UserService {
	return &UserService{
		otp: NewOTPManager(),
	}
}

func (us *UserService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	return app.GetUserByID(userID)
}

func (us *UserService) EditUser(userID uuid.UUID, updatedUser *models.User) error {
	existingUser, err := app.GetUserByID(userID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	if updatedUser.Username != "" {
		existingUser.Username = updatedUser.Username
	}

	if updatedUser.Password != "" {
		existingUser.Password = updatedUser.Password
	}

	if updatedUser.FirstName != "" {
		existingUser.FirstName = updatedUser.FirstName
	}

	if updatedUser.MiddleName != nil {
		existingUser.MiddleName = updatedUser.MiddleName
	}

	if updatedUser.LastName != "" {
		existingUser.LastName = updatedUser.LastName
	}

	if updatedUser.Email != "" {
		existingUser.Email = updatedUser.Email
	}

	if updatedUser.StudentID != "" {
		existingUser.StudentID = updatedUser.StudentID
	}

	if updatedUser.Major != "" {
		existingUser.Major = updatedUser.Major
	}

	if updatedUser.Year != "" {
		existingUser.Year = updatedUser.Year
	}

	if updatedUser.DateOfBirth != nil {
		existingUser.DateOfBirth = updatedUser.DateOfBirth
	}

	if updatedUser.InstitutionName != nil {
		existingUser.InstitutionName = updatedUser.InstitutionName
	}

	if updatedUser.Gender != "" {
		existingUser.Gender = updatedUser.Gender
	}

	return app.UpdateUser(userID, existingUser)
}

func (us *UserService) DeleteUser(userID uuid.UUID) error {
	return app.DeleteUser(userID)
}

func (us *UserService) ChangePassword(userID uuid.UUID, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return app.ChangePassword(userID, string(hashedPassword))
}

func (us *UserService) GetUserByUsername(username string) (*models.User, error) {
	return app.GetUserByUsername(username)
}

func (us *UserService) GetUserByEmail(email string) (*models.User, error) {
	return app.GetUserByEmail(email)
}

func (us *UserService) GetRoleIDByUserID(userID uuid.UUID) (int, error) {
	return app.GetRoleIDByUserID(userID)
}

func (us *UserService) CheckStudentIDExists(studentID string) (bool, error) {
	return app.CheckStudentIDExists(studentID)
}

func (us *UserService) ListUsers() ([]models.User, error) {
	log.Println("service list users")
	return app.ListUsers()
}

func (us *UserService) AdminUpdateRoleAndStudentIDVerified(userID uuid.UUID, roleID int, studentIDVerified bool) error {
	return app.AdminUpdateRoleAndStudentIDVerified(userID, roleID, studentIDVerified)
}

func (us *UserService) UploadProfilePicture(userID uuid.UUID, profilePicture string) error {
	return app.UploadProfilePicture(userID, profilePicture)
}

func (us *UserService) UploadStudentID(userID uuid.UUID, profilePicture string) error {
	return app.UploadStudentID(userID, profilePicture)
}

func (us *UserService) EnableTwoFA(userID uuid.UUID) (string, string, error) {
	user, err := app.GetUserByID(userID)
	if err != nil {
		return "", "", err
	}

	log.Println("Generating TOTP key")

	secret, err := utils.GenerateTOTPKey(user.Email)
	if err != nil {
		return "", "", err
	}

	log.Println("Generated TOTP secret:", secret)

	qr, err := utils.GenerateQRCodeBase64(secret)
	if err != nil {
		return "", "", err
	}

	secretStr := secret.Secret()

	user.TwoFASecret = &secretStr
	user.TwoFAImage = &qr

	err = app.SaveTwoFAInfo(userID, secretStr, qr)
	if err != nil {
		return "", "", err
	}

	return qr, secretStr, nil
}

func (us *UserService) VerifyTwoFA(userID uuid.UUID, code string) (bool, error) {
	user, err := app.GetUserByID(userID)
	if err != nil {
		return false, err
	}

	if user.TwoFASecret == nil {
		log.Println("No TOTP secret found for user")
		return false, fmt.Errorf("no TOTP secret found for user")
	}

	log.Println("Verifying TOTP code with secret:", *user.TwoFASecret)
	log.Println("TOTP code to verify:", code)

	// Ensure the correct settings for TOTP validation
	valid, err := totp.ValidateCustom(code, *user.TwoFASecret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA256,
	})

	if err != nil {
		return false, err
	}

	log.Println("Is TOTP code valid?", valid)
	return valid, nil
}

func (us *UserService) ChangeTwoFAStatus(userID uuid.UUID, enable bool) error {
	user, err := app.GetUserByID(userID)
	if err != nil {
		return err
	}

	user.TwoFAEnabled = enable

	err = app.ToggleTwoFA(userID, enable)
	if err != nil {
		return err
	}

	return nil
}
