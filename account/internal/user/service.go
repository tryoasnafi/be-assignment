package user

import (
	"sync"

	"github.com/google/uuid"
)

var (
	service *UserService
	once sync.Once
)

type Service interface {
	GetByID(userUUID uuid.UUID) (*User, error)
	Create(user *User) (*User, error)
}

type UserService struct {
	repo Repository
}

func NewService(repo Repository) *UserService {
	once.Do(func ()  {
		service = &UserService{repo}
	})
	return service
}

func (s UserService) GetByID(userUUID uuid.UUID) (*User, error) {
	return s.repo.GetByID(userUUID)
}

func (s UserService) Create(user *User) (*User, error) {
	return s.repo.Create(user)
}
