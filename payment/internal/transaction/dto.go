package transaction

import (
	"time"

	"github.com/google/uuid"
)

type TransactionOperation string

const (
	OperationCredit TransactionOperation = "credit"
	OperationDebit  TransactionOperation = "debit"
)

type TransactionType string

const (
	TransactionWithdraw TransactionType = "withdraw"
	TransactionDeposit  TransactionType = "deposit"
	TransactionSend     TransactionType = "send"
)

type Transaction struct {
	ID                 uuid.UUID       `json:"id" gorm:"type:uuid;default:gen_random_uuid()"`
	AccountId          uint            `json:"account_id"`
	RecipientAccountId uint            `json:"recipient_account_id,omitempty"`
	Type               TransactionType `json:"type"`
	Amount             float64         `json:"amount"`
	Currency           string          `json:"currency"`
	Status             bool            `json:"status"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

type SendRequest struct {
	AccountID          uint    `json:"account_id"`
	RecipientAccountID uint    `json:"recipient_account_id"`
	Amount             float64 `json:"amount"`
	Currency           string  `json:"currency"`
}

type SendResponse struct {
	SendRequest
	TransactionID uuid.UUID `json:"transaction_id"`
	TransactionAt time.Time `json:"transaction_at"`
}

type WithdrawRequest struct {
	AccountID uint    `json:"account_id"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
}

type WithdrawResponse struct {
	WithdrawRequest
	TransactionID uuid.UUID `json:"transaction_id"`
	TransactionAt time.Time `json:"transaction_at"`
}
