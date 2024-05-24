package services

import (
	"Backend/configs"
	"Backend/internal/database/app"
	"Backend/internal/models"
	"Backend/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"net/http"
	"time"
)

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (as *AuthService) RegisterUser(user *models.User) error {
	// Email validation
	if err := as.ValidateEmail(user.Email); err != nil {
		return err
	}

	user.ID = uuid.New()
	user.RoleID = 2

	// Set major based on studentID
	if user.StudentID[:3] == "001" {
		user.Major = "informatics"
	} else if user.StudentID[:3] == "012" {
		user.Major = "information system"
	} else if user.StudentID[:3] == "013" {
		user.Major = "visual communication design"
	} else if user.StudentID[:3] == "025" {
		user.Major = "interior design"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	err = app.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (as *AuthService) LoginUser(usernameOrEmail string, password string) (*models.User, error) {
	user, err := app.AuthenticateUser(usernameOrEmail)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, &utils.UnauthorizedError{Message: "invalid credentials"}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Error comparing hashed password: %v", err)
		return nil, &utils.UnauthorizedError{Message: "invalid credentials"}
	}

	return user, nil
}

func (as *AuthService) IsUsernameExists(username string) (bool, error) {
	return app.IsUsernameExists(username)
}

func (as *AuthService) IsEmailExists(email string) (bool, error) {
	return app.IsEmailExists(email)
}

func (as *AuthService) IsStudentIDExists(studentID string) (bool, error) {
	return app.CheckStudentIDExists(studentID)
}

func (as *AuthService) GetUserByStudentID(studentID string) (*models.User, error) {
	return app.GetUserByStudentID(studentID)
}

func (as *AuthService) GetUserByUsernameOrEmail(usernameOrEmail string) (*models.User, error) {
	return app.AuthenticateUser(usernameOrEmail)
}

func (as *AuthService) CheckStudentIDExists(studentID string) (bool, error) {
	return app.CheckStudentIDExists(studentID)
}

func (as *AuthService) IsEmailVerified(username string) (bool, error) {
	return app.IsEmailVerified(username)
}

func (as *AuthService) IsTokenVerificationEmailExists(token string) (bool, error) {
	return app.IsTokenVerificationEmailExists(token)
}

func (as *AuthService) UpdateEmailVerificationToken(email, token string) error {
	return app.UpdateEmailVerificationToken(email, token)
}

func (as *AuthService) VerifyEmail(token string) error {
	return app.VerifyEmail(token)
}

type HunterEmailVerification struct {
	Data struct {
		Status     string `json:"status"`
		Result     string `json:"result"`
		Score      int    `json:"score"`
		Regexp     bool   `json:"regexp"`
		Gibberish  bool   `json:"gibberish"`
		Disposable bool   `json:"disposable"`
		Webmail    bool   `json:"webmail"`
		MxRecords  bool   `json:"mx_records"`
		SmtpServer bool   `json:"smtp_server"`
		SmtpCheck  bool   `json:"smtp_check"`
		AcceptAll  bool   `json:"accept_all"`
		Block      bool   `json:"block"`
	} `json:"data"`
}

func (as *AuthService) ValidateEmail(email string) error {
	// Load the Hunter API key from the config
	apiKey := configs.LoadConfig().HunterApiKey
	url := fmt.Sprintf("https://api.hunter.io/v2/email-verifier?email=%s&api_key=%s", email, apiKey)

	// Create a new HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(resp.Body)

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to validate email: received status code %d", resp.StatusCode)
	}

	// Parse the JSON response
	var verification HunterEmailVerification
	if err := json.NewDecoder(resp.Body).Decode(&verification); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Check the email status
	switch verification.Data.Status {
	case "valid":
		// Email is valid, proceed with registration
		return nil
	case "invalid":
		return errors.New("the email address is invalid")
	case "disposable":
		return errors.New("the email address is from a disposable email service")
	case "webmail":
		// Optionally handle webmail addresses differently
		return nil
	case "unknown":
		return errors.New("failed to verify the email address")
	default:
		return errors.New("unexpected status from email verification")
	}
}
