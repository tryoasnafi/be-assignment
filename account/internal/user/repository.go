package user

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface  {
	GetByID(userUUID uuid.UUID) (*User, error)
	// GetDeletedByID(id uint) (*User, error)
	Create(user *User) (*User, error)
}

type UserRepository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (repo *UserRepository) GetByID(userUUID uuid.UUID) (*User, error) {
	user := &User{}
	result := repo.DB.First(user, "uuid = ?", userUUID)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (repo *UserRepository) Create(user *User) (*User, error) {
	result := repo.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("user %v already exists", user)
	}
	return user, nil
}
