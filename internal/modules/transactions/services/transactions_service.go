package transactionsservices

import (
	"context"
	"errors"
	"fmt"
	"my_finance/internal/apperrors"
	"my_finance/internal/constants"
	loggerHelper "my_finance/internal/logger"
	categories_model "my_finance/models/categories"
	transactions_model "my_finance/models/transactions"
	"time"
)

type TransactionsRepository interface {
	GetAll(r context.Context) ([]transactions_model.TransactionsModel, error)
	Create(r context.Context, transaction transactions_model.TransactionsModel) (int, error)
	FindById(r context.Context, id int) (transactions_model.TransactionsModel, error)
	Pay(r context.Context, id int, paidAt time.Time) error
}

type CategoryRepository interface {
	FindById(r context.Context, id int) (categories_model.CategoryModel, error)
}

type TransactionsService struct {
	repository         TransactionsRepository
	categoryRepository CategoryRepository
}

func NewTransactionsService(repository TransactionsRepository, categoryRepository CategoryRepository) *TransactionsService {
	return &TransactionsService{
		repository:         repository,
		categoryRepository: categoryRepository,
	}
}

func (s *TransactionsService) GetAll(ctx context.Context) ([]transactions_model.TransactionsModel, error) {
	transactions, err := s.repository.GetAll(ctx)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao executar o select:", err)

		return []transactions_model.TransactionsModel{}, err
	}

	return transactions, nil
}

func (s *TransactionsService) Create(ctx context.Context, transaction transactions_model.TransactionsModel) (int, error) {
	if err := transaction.Validate(); err != nil {
		return 0, err
	}

	_, err := s.categoryRepository.FindById(ctx, transaction.CategoryId)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return 0, fmt.Errorf("Categoria %d não localizada.", transaction.CategoryId)
		}

		loggerHelper.ErrorLogger.Println("Erro ao localizar a categoria da transação:", err)
		return 0, err
	}

	transactionId, err := s.repository.Create(ctx, transaction)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao cadastrar a transação:", err.Error())
		return 0, err
	}

	return transactionId, nil
}

func (t *TransactionsService) FindById(ctx context.Context, id int) (transactions_model.TransactionsModel, error) {
	if id <= 0 {
		return transactions_model.TransactionsModel{}, fmt.Errorf("O ID da transação precisa ser maior que zero.")
	}

	transaction, err := t.repository.FindById(ctx, id)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao localizar a transação:", err.Error())
		return transactions_model.TransactionsModel{}, err
	}

	return transaction, nil
}

func (t *TransactionsService) Pay(ctx context.Context, id int, paidAt time.Time) error {
	if id <= 0 {
		return fmt.Errorf("O ID da transação precisa ser maior que zero.")
	}

	transaction, err := t.FindById(ctx, id)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao localizar a transação:", err)
		return fmt.Errorf("A transação %d não foi localizada.", id)
	}

	if transaction.Status == constants.StatusCancelado {
		return fmt.Errorf("A transação %d está cancelada.", id)
	}

	if transaction.Status == constants.StatusPago {
		return fmt.Errorf("A transação %d está paga.", id)
	}

	if transaction.PaidAt != nil {
		return fmt.Errorf("A transação %d já possui uma data de pagamento.", id)
	}

	if transaction.DeletedAt != nil {
		return fmt.Errorf("A transação %d está deletada.", id)
	}

	if err := t.repository.Pay(ctx, id, paidAt); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao pagar a transação:", err)
		return err
	}

	return nil
}
