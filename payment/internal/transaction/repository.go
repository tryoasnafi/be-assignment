package transaction

import (
	"fmt"
	"log"

	. "github.com/tryoasnafi/be-assignment/payment/internal/transaction/dto"
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
		AccountId:          sendReq.AccountID,
		RecipientAccountId: sendReq.RecipientAccountID,
		Amount:             sendReq.Amount,
		Type:               TransactionSend,
		Currency:           sendReq.Currency,
		Status:             true,
	}

	senderAccount := Account{}
	result := repo.DB.First(&senderAccount, sendReq.AccountID)
	if result.Error != nil {
		return Transaction{}, result.Error
	}
	if senderAccount.Balance < sendReq.Amount {
		return Transaction{}, fmt.Errorf("insufficient balance")
	}
	sender := AccountHistories{
		AccountID:    sendReq.AccountID,
		Type:         OperationDebit,
		Amount:       sendReq.Amount,
		BeforeAmount: senderAccount.Balance,
		FinalAmount:  senderAccount.Balance - sendReq.Amount,
	}

	recipientAccount := Account{}
	result = repo.DB.First(&recipientAccount, sendReq.RecipientAccountID)
	if result.Error != nil {
		return Transaction{}, result.Error
	}
	recipient := AccountHistories{
		AccountID:    sendReq.RecipientAccountID,
		Type:         OperationCredit,
		Amount:       sendReq.Amount,
		BeforeAmount: recipientAccount.Balance,
		FinalAmount:  recipientAccount.Balance + sendReq.Amount,
	}
	err := repo.DB.Transaction(func(tx *gorm.DB) error {
		tx.Create(&transaction)
		// update accounts histories
		sender.TransactionID = transaction.ID
		recipient.TransactionID = transaction.ID
		tx.Create(&sender)
		tx.Create(&recipient)
		// update accounts balance
		tx.Model(&senderAccount).Update("balance", sender.FinalAmount)
		tx.Model(&recipientAccount).Update("balance", recipient.FinalAmount)
		return nil
	})
	if err != nil {
		return Transaction{}, err
	}
	return transaction, nil
}

func (repo transactionRepository) Withdraw(withdrawReq WithdrawRequest) (Transaction, error) {
	transaction := Transaction{
		AccountId: withdrawReq.AccountID,
		Amount:    withdrawReq.Amount,
		Type:      TransactionWithdraw,
		Currency:  withdrawReq.Currency,
		Status:    true,
	}
	account := Account{}
	result := repo.DB.First(&account, withdrawReq.AccountID)
	if result.Error != nil {
		log.Println("account", account, result.Error)
		return Transaction{}, result.Error
	}
	if account.Balance < withdrawReq.Amount {
		return Transaction{}, fmt.Errorf("insufficient balance")
	}
	accountHistory := AccountHistories{
		AccountID:    withdrawReq.AccountID,
		Type:         OperationDebit,
		Amount:       withdrawReq.Amount,
		BeforeAmount: account.Balance,
		FinalAmount:  account.Balance - withdrawReq.Amount,
	}
	err := repo.DB.Transaction(func(tx *gorm.DB) error {
		tx.Create(&transaction)
		// update account histories sender and recipient
		accountHistory.TransactionID = transaction.ID
		tx.Create(&accountHistory)
		// update account balance
		tx.Model(&account).Update("balance", accountHistory.FinalAmount)
		return nil
	})
	if err != nil {
		return Transaction{}, err
	}
	return transaction, nil
}
