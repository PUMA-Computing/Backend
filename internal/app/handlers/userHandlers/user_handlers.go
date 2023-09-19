package userHandlers

import (
	"Backend/internal/app/domain/user"
	user2 "Backend/internal/app/interfaces/service/userService"
	"Backend/internal/utils/token"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

type UserHandlers struct {
	userService user2.UserServices
}

func NewUserHandlers(userService user2.UserServices) *UserHandlers {
	return &UserHandlers{userService: userService}
}

func (h *UserHandlers) RegisterUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user *user.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error parsing registration data", "error": err.Error()})
		}

		if err := h.validateRegistrationData(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		if err := h.userService.RegisterUser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error registering user", "error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
	}
}

func (h *UserHandlers) validateRegistrationData(user *user.User) error {
	if user.Email == "" {
		return errors.New("email is required")
	}

	if user.Password == "" {
		return errors.New("password is required")
	}

	if user.FirstName == "" {
		return errors.New("first name is required")
	}

	if user.LastName == "" {
		return errors.New("last name is required")
	}

	if user.NIM == "" {
		return errors.New("nim is required")
	}

	if user.Year == "" {
		return errors.New("year is required")
	}

	return nil
}

func (h *UserHandlers) Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var loginData struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.BodyParser(&loginData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error parsing login data"})
		}

		user, err := h.userService.AuthenticateUser(loginData.Email, loginData.Password)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid email or password"})
		}

		sessionToken, err := token.GenerateJWTToken(user.User.ID, user.User.RoleID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error generating session token", "error": err.Error()})
		}

		fmt.Printf("Generated Token Claims:", sessionToken)

		expirationTime := time.Now().Add(token.SessionDuration)
		if err := token.StoreSessionData(user.User.ID, sessionToken, expirationTime); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error storing session data", "error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User logged in successfully", "access_token": sessionToken})
	}
}

func (h *UserHandlers) Logout() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userIDInterface, ok := c.Locals("userID").(uuid.UUID)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user id format"})
		}

		if err := h.userService.Logout(userIDInterface); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error logout user", "error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User logged out successfully"})
	}
}

func (h *UserHandlers) GetUserProfile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(uuid.UUID)
		user, err := h.userService.GetUserByID(userID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error getting user profile"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User profile retrieved successfully", "user": user})
	}
}

func (h *UserHandlers) GetAllUsers() fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := h.userService.GetAllUsers()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error getting all users"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Users retrieved successfully", "users": users})
	}
}

func (h *UserHandlers) GetUserByEmail() fiber.Handler {
	return func(c *fiber.Ctx) error {
		email := c.Locals("email").(string)
		user, err := h.userService.GetUserByEmail(email)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error getting user by email"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User retrieved successfully", "user": user})
	}
}

func (h *UserHandlers) GetUserRoleByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(uuid.UUID)
		userRoleID, err := h.userService.GetUserRoleByID(userID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error getting user role by id"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User role retrieved successfully", "userRoleID": userRoleID})
	}
}

func (h *UserHandlers) GetUserRoleByEmail() fiber.Handler {
	return func(c *fiber.Ctx) error {
		email := c.Locals("email").(string)
		userRoleID, err := h.userService.GetUserRoleByEmail(email)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error getting user role by email"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User role retrieved successfully", "userRoleID": userRoleID})
	}
}

func (h *UserHandlers) UpdateUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user *user.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error parsing user data"})
		}

		if err := h.userService.UpdateUser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error updating user"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User updated successfully"})
	}
}

func (h *UserHandlers) DeleteUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(uuid.UUID)
		if err := h.userService.DeleteUser(userID); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error deleting user"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User deleted successfully"})
	}
}
