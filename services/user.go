package services

import (
	"errors"

	"github.com/wanliqun/go-wallet-app/models"
	"gorm.io/gorm"
)

var (
	ErrUnauthorized = errors.New("Unauthorized")
	ErrUserNotFound = errors.New("user not found")

	_ IUserService = &UserService{}
)

type IUserService interface {
	GetUserByName(name string) (*models.User, bool, error)
}

// UserService represents the service for user-related operations
type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (svc *UserService) GetUserByName(name string) (*models.User, bool, error) {
	var user models.User
	if err := svc.DB.Where("name = ?", name).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &user, true, nil
}
