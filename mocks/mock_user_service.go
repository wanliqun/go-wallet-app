package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/wanliqun/go-wallet-app/models"
	"github.com/wanliqun/go-wallet-app/services"
)

var (
	_ services.IUserService = &MockUserService{}
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserByName(name string) (*models.User, bool, error) {
	args := m.Called(name)
	return args.Get(0).(*models.User), args.Bool(1), args.Error(2)
}
