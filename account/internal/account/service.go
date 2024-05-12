package account

import (
	"time"

	"github.com/google/uuid"

	. "github.com/tryoasnafi/be-assignment/common/model"
)

type Service interface {
	GetAllAccounts(userID uuid.UUID) ([]Account, error)
	GetAccountByID(id uint) (*Account, error)
	CreateAccount(userUUID uuid.UUID, accountReq CreateAccountRequest) (*Account, error)
}

type accountService struct {
	repo Repository
}

func NewService(repo Repository) *accountService {
	return &accountService{repo}
}

func (s accountService) GetAllAccounts(userID uuid.UUID) ([]Account, error) {
	return s.repo.GetAccounts(userID)
}

func (s accountService) GetAccountByID(id uint) (*Account, error) {
	return s.repo.Find(id)
}

func (s accountService) CreateAccount(userUUID uuid.UUID, accountReq CreateAccountRequest) (*Account, error) {
	account := &Account{
		UserID:   userUUID,
		Type:     accountReq.Type,
		Currency: accountReq.Currency,
		Balance:  0,
		Status:   AccountActive,
		OpenedAt: time.Now(),
	}
	return s.repo.Create(account)
}
