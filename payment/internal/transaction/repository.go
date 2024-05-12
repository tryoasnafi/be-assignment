package transaction

import (
	"gorm.io/gorm"
)

type Repository interface {
	Send(sendReq SendRequest) (Transaction, error)
	Withdraw(withdrawReq WithdrawRequest) (Transaction, error)
}

type transactionRepository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) transactionRepository {
	return transactionRepository{db}
}

func (repo transactionRepository) Send(sendReq SendRequest) (Transaction, error) {
	transaction := Transaction{
		AccountId: sendReq.AccountID,
		RecipientAccountId: sendReq.RecipientAccountID,
		Amount: sendReq.Amount,
		Type: TransactionSend,
		Currency: sendReq.Currency,
		Status: true,
	}
	err := repo.DB.Transaction(func(tx *gorm.DB) error {
		tx.Create(&transaction);
		return nil
	});
	if err != nil {
		return Transaction{}, err
	}
	return transaction, nil
}

func (repo transactionRepository) Withdraw(withdrawReq WithdrawRequest) (Transaction, error) {
	transaction := Transaction{
		AccountId: withdrawReq.AccountID,
		Amount: withdrawReq.Amount,
		Type: TransactionWithdraw,
		Currency: withdrawReq.Currency,
		Status: true,
	}
	err := repo.DB.Transaction(func(tx *gorm.DB) error {
		tx.Create(&transaction);
		return nil
	});
	if err != nil {
		return Transaction{}, err
	}
	return transaction, nil
}

