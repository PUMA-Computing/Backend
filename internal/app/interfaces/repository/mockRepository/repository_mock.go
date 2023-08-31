package mockRepository

import (
	"Backend/internal/app/domain/user"
	user2 "Backend/internal/app/interfaces/service/userService"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func NewMockUserService() *MockUserService {
	return &MockUserService{}
}

func (m *MockUserService) RegisterUser(user *user.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) AuthenticateUser(email, password string) (*user2.AuthResponse, error) {
	args := m.Called(email, password)
	return args.Get(0).(*user2.AuthResponse), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(email string) (*user.User, error) {
	args := m.Called(email)
	return args.Get(0).(*user.User), args.Error(1)
}
