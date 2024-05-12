package transaction

import (
	. "github.com/tryoasnafi/be-assignment/payment/internal/transaction/dto"
)

type Service interface {
	Send(req SendRequest) (SendResponse, error)
	Withdraw(req WithdrawRequest) (WithdrawResponse, error)
}

type transactionService struct {
	repo Repository
}

func NewService(repo Repository) transactionService {
	return transactionService{repo}
}

func (s transactionService) Send(req SendRequest) (SendResponse, error) {
	tx, err := s.repo.Send(req)
	if err != nil {
		return SendResponse{}, err
	}
	return SendResponse{
		SendRequest:   req,
		TransactionID: tx.ID,
		TransactionAt: tx.CreatedAt,
	}, nil
}

func (s transactionService) Withdraw(req WithdrawRequest) (WithdrawResponse, error) {
	tx, err := s.repo.Withdraw(req)
	if err != nil {
		return WithdrawResponse{}, err
	}
	return WithdrawResponse{
		WithdrawRequest: req,
		TransactionID:   tx.ID,
		TransactionAt:   tx.CreatedAt,
	}, nil
}
