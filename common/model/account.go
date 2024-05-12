package model

import (
	"time"

	"github.com/google/uuid"
)

type AccountType string

const (
	AccountSavings  AccountType = "savings"
	AccountChecking AccountType = "checking"
)

type AccountStatus string

const (
	AccountActive    AccountStatus = "active"
	AccountSuspended AccountStatus = "suspended"
	AccountClosed    AccountStatus = "closed"
)


type Account struct {
	ID        uint                `json:"id"`
	UserID    uuid.UUID           `json:"user_id" gorm:"type:uuid"`
	Type      AccountType         `json:"type"`
	Currency  string              `json:"currency"`
	Balance   float64             `json:"balance"`
	Status    AccountStatus       `json:"status"`
	OpenedAt  time.Time           `json:"opened_at"`
	ClosedAt  *time.Time          `json:"closed_at,omitempty"`
	Histories []*AccountHistories `json:"histories"`
}

type AccountHistories struct {
	ID            uint                 `json:"id"`
	AccountID     uint                 `json:"account_id"`
	TransactionID uuid.UUID            `json:"transaction_id" gorm:"type:uuid"`
	Type          TransactionOperation `json:"type"`
	Amount        float64              `json:"amount"`
	BeforeAmount  float64              `json:"before_amount"`
	FinalAmount   float64              `json:"final_amount"`
	CreatedAt     time.Time            `json:"created_at"`
}