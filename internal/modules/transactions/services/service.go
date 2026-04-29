package transactionsservices

import (
	"context"
	loggerHelper "my_finance/internal/logger"
	transactions_model "my_finance/models/transactions"
)

type TransactionsRepository interface {
	GetAll(r context.Context) ([]transactions_model.TransactionsModel, error)
	Create(r context.Context, transaction transactions_model.TransactionsModel) (int, error)
	FindById(r context.Context, id int) (transactions_model.TransactionsModel, error)
}

type TransactionsService struct {
	repository TransactionsRepository
}

func NewTransactionsService(repository TransactionsRepository) *TransactionsService {
	return &TransactionsService{
		repository: repository,
	}
}

func (s *TransactionsService) GetAll(ctx context.Context) ([]transactions_model.TransactionsModel, error) {
	transactions, err := s.repository.GetAll(ctx)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao retornar as transações:", err)

		return []transactions_model.TransactionsModel{}, err
	}

	return transactions, nil
}

func (s *TransactionsService) Create(ctx context.Context, transaction transactions_model.TransactionsModel) (int, error) {
	if err := transaction.Validate(); err != nil {
		return 0, err
	}

	transactionId, err := s.repository.Create(ctx, transaction)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao cadastrar a transição:", err.Error())
		return 0, err
	}

	return transactionId, nil
}

func (t *TransactionsService) FindById(ctx context.Context, id int) (transactions_model.TransactionsModel, error) {
	return transactions_model.TransactionsModel{}, nil
}
