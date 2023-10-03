package user

import (
	"Backend/internal/services"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	router := gin.New()
	userService := &services.UserService{}
	permissionService := &services.PermissionService{}
	handlers := NewUserHandlers(userService, permissionService)
	router.POST("/register", handlers.RegisterUser)

	testPayload := []byte(`{
	"username": "testuser",
	"password": "testpassword",
	"email": "testmail@gmail.com",
	"first_name": "John",
	"last_name": "Doe",
	"student_id": "00120220001",
	"Major": "Information Technology",
	"RoleID": 1
    }`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(testPayload))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Data struct {
			Username string `json:"username"`
		} `json:"data"`
		Meta struct {
			Message string `json:"message"`
		} `json:"meta"`
	}
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", response.Data.Username)
	assert.Equal(t, "User Registered Successfully", response.Meta.Message)
}
