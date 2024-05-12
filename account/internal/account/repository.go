package account

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetAccounts(userID uuid.UUID) ([]Account, error)
	Find(id uint) (*Account, error)
	Create(acc *Account) (*Account, error)
	// GetTransactions(id uint) ([]Transaction, error)
}

type AccountRepository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) AccountRepository {
	return AccountRepository{db}
}

func (repo AccountRepository) GetAccounts(userID uuid.UUID) ([]Account, error) {
	var accounts []Account
	result := repo.DB.Find(&accounts, "user_id = ?", userID)
	return accounts, result.Error
}

func (repo AccountRepository) Find(id uint) (*Account, error) {
	account := &Account{}
	result := repo.DB.First(account, id)
	return account, result.Error
}

func (repo AccountRepository) Create(account *Account) (*Account, error) {
	result := repo.DB.Create(account)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("account %v already exists", account)
	}
	return account, nil
}
