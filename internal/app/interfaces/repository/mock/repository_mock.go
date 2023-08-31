package mock

import (
	"Backend/internal/app/domain/user"
	user2 "Backend/internal/app/interfaces/service/user"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func NewMockUserService() *MockUserService {
	return &MockUserService{}
}

func (m *MockUserService) RegisterUser(user *user.User) error {
	ret := m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*user.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (m *MockUserService) AuthenticateUser(email, password string) (*user2.AuthResponse, error) {
	args := m.Called(email, password)
	return args.Get(0).(*user2.AuthResponse), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(email string) (*user.User, error) {
	args := m.Called(email)
	return args.Get(0).(*user.User), args.Error(1)
}
