package account

import (
	"time"

	"github.com/google/uuid"
	. "github.com/tryoasnafi/be-assignment/common/model"
)

type CreateAccountRequest struct {
	UserID   uuid.UUID       `json:"user_id"`
	Type     AccountType `json:"type"`
	Currency string          `json:"currency"`
}

type AccountHistoriesRequest struct {
	AccountID uint      `json:"account_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}
