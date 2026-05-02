package transactionsservices

import (
	"context"
	"errors"
	"fmt"
	"my_finance/internal/apperrors"
	transactionsconstants "my_finance/internal/constants/transactions"
	loggerHelper "my_finance/internal/logger"
	categories_model "my_finance/models/categories"
	incomereceiptsmodel "my_finance/models/income_receipts"
	transactions_model "my_finance/models/transactions"
)

type TransactionsRepository interface {
	GetAll(r context.Context) ([]transactions_model.TransactionsModel, error)
	Create(r context.Context, transaction transactions_model.TransactionsModel) (transactions_model.TransactionsModel, error)
	Update(r context.Context, transaction transactions_model.TransactionsModel) (transactions_model.TransactionsModel, error)
	FindById(r context.Context, id int) (transactions_model.TransactionsModel, error)
	Pay(r context.Context, id int) error
	Cancel(r context.Context, id int) error
	Delete(r context.Context, id int) error
}

type CategoryRepository interface {
	FindById(r context.Context, id int) (categories_model.CategoryModel, error)
}

type IncomeReceiptsRepository interface {
	FindById(ctx context.Context, id int) (*incomereceiptsmodel.IncomeReceiptsModel, error)
	Create(ctx context.Context, incomeReceipt incomereceiptsmodel.IncomeReceiptsModel) (incomereceiptsmodel.IncomeReceiptsModel, error)
}

type TransactionsService struct {
	repository               TransactionsRepository
	categoryRepository       CategoryRepository
	incomeReceiptsRepository IncomeReceiptsRepository
}

func NewTransactionsService(
	repository TransactionsRepository,
	categoryRepository CategoryRepository,
	incomeReceiptsRepository IncomeReceiptsRepository,
) *TransactionsService {
	return &TransactionsService{
		repository:               repository,
		categoryRepository:       categoryRepository,
		incomeReceiptsRepository: incomeReceiptsRepository,
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

func (s *TransactionsService) Create(ctx context.Context, transaction transactions_model.TransactionsModel) (transactions_model.TransactionsModel, error) {
	if err := transaction.Validate(); err != nil {
		return transactions_model.TransactionsModel{}, err
	}

	if transaction.Origin != "" && transaction.Origin != transactionsconstants.ManualOrigin {
		loggerHelper.InfoLogger.Println("Informações de recebimento informados em uma transação manual, vai converter para criação manual")

		transaction.Origin = transactionsconstants.ManualOrigin
		transaction.OriginId = nil
	}

	_, err := s.categoryRepository.FindById(ctx, transaction.CategoryId)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			loggerHelper.ErrorLogger.Printf("Categoria %d não localizada.", transaction.CategoryId)

			return transactions_model.TransactionsModel{}, apperrors.NewErrNotFound(fmt.Sprintf("Categoria %d não localizada.", transaction.CategoryId))
		}

		loggerHelper.ErrorLogger.Println("Erro ao localizar a categoria da transação:", err)
		return transactions_model.TransactionsModel{}, err
	}

	if transaction.PaidAt != nil {
		transaction.Status = transactionsconstants.StatusPago

	}

	transactionData, err := s.repository.Create(ctx, transaction)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao criar a transação:", err.Error())
		return transactions_model.TransactionsModel{}, err
	}

	return transactionData, nil
}

func (s *TransactionsService) Update(ctx context.Context, transaction transactions_model.TransactionsModel) (transactions_model.TransactionsModel, error) {
	if err := transaction.Validate(); err != nil {
		return transactions_model.TransactionsModel{}, err
	}

	if transaction.Origin != "" && transaction.Origin != transactionsconstants.ManualOrigin {
		loggerHelper.InfoLogger.Println("Informações de recebimento informados em uma transação manual, vai converter para criação manual")

		transaction.Origin = transactionsconstants.ManualOrigin
		transaction.OriginId = nil
	}

	_, err := s.categoryRepository.FindById(ctx, transaction.CategoryId)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			loggerHelper.ErrorLogger.Printf("Categoria %d não localizada.", transaction.CategoryId)

			return transactions_model.TransactionsModel{}, apperrors.NewErrNotFound(fmt.Sprintf("Categoria %d não localizada.", transaction.CategoryId))
		}

		loggerHelper.ErrorLogger.Println("Erro ao localizar a categoria da transação:", err)
		return transactions_model.TransactionsModel{}, err
	}

	transactionData, err := s.FindById(ctx, transaction.Id)

	if transactionData.DeletedAt != nil {
		return transactions_model.TransactionsModel{}, fmt.Errorf("Transação %d já deletada.", transactionData.Id)
	}

	if transactionData.Status == transactionsconstants.StatusCancelado {
		return transactions_model.TransactionsModel{}, fmt.Errorf("Transação %d já cancelada.", transactionData.Id)
	}

	if transactionData.Status == transactionsconstants.StatusPago {
		return transactions_model.TransactionsModel{}, fmt.Errorf("Transação %d já paga.", transactionData.Id)
	}

	transactionData, err = s.repository.Update(ctx, transaction)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao alterar a transação:", err.Error())
		return transactions_model.TransactionsModel{}, err
	}

	return transactionData, nil
}

func (s *TransactionsService) FindById(ctx context.Context, id int) (transactions_model.TransactionsModel, error) {
	if id <= 0 {
		return transactions_model.TransactionsModel{}, fmt.Errorf("O ID da transação precisa ser maior que zero.")
	}

	transaction, err := s.repository.FindById(ctx, id)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			loggerHelper.ErrorLogger.Printf("Transação %d não localizada.", id)
			return transactions_model.TransactionsModel{}, apperrors.NewErrNotFound(fmt.Sprintf("Transação %d não localizada.", id))
		}

		loggerHelper.ErrorLogger.Println("Erro ao localizar a transação:", err.Error())
		return transactions_model.TransactionsModel{}, err
	}

	return transaction, nil
}

func (s *TransactionsService) Pay(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("O ID da transação precisa ser maior que zero.")
	}

	transaction, err := s.FindById(ctx, id)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			loggerHelper.ErrorLogger.Printf("Transação %d não localizada.", id)

			return apperrors.NewErrNotFound(fmt.Sprintf("Transação %d não localizada.", id))
		}

		loggerHelper.ErrorLogger.Println("Erro ao localizar a transação:", err.Error())
		return err
	}

	if transaction.DeletedAt != nil {
		return fmt.Errorf("Transação %d já deletada.", id)
	}

	if transaction.Status == transactionsconstants.StatusCancelado {
		return fmt.Errorf("Transação %d já cancelada.", id)
	}

	if transaction.Status == transactionsconstants.StatusPago {
		return fmt.Errorf("Transação %d já paga.", id)
	}

	if transaction.PaidAt != nil {
		return fmt.Errorf("Transação %d já possui uma data de pagamento.", id)
	}

	if err := s.repository.Pay(ctx, id); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao pagar a transação:", err)
		return err
	}

	return nil
}

func (s *TransactionsService) Cancel(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("O ID da transação precisa ser maior que zero.")
	}

	transaction, err := s.FindById(ctx, id)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			loggerHelper.ErrorLogger.Printf("Transação %d não localizada.", id)

			return apperrors.NewErrNotFound(fmt.Sprintf("Transação %d não localizada.", id))
		}

		loggerHelper.ErrorLogger.Println("Erro ao localizar a transação:", err.Error())
		return err
	}

	if transaction.DeletedAt != nil {
		return fmt.Errorf("Transação %d já deletada.", id)
	}

	if transaction.Status == transactionsconstants.StatusCancelado {
		return fmt.Errorf("Transação %d já cancelada.", id)
	}

	if err := s.repository.Cancel(ctx, id); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao pagar a transação:", err)
		return err
	}

	return nil
}

func (s *TransactionsService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("O ID da transação precisa ser maior que zero.")
	}

	transaction, err := s.FindById(ctx, id)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			loggerHelper.ErrorLogger.Printf("Transação %d não localizada.", id)

			return apperrors.NewErrNotFound(fmt.Sprintf("Transação %d não localizada.", id))
		}

		loggerHelper.ErrorLogger.Println("Erro ao localizar a transação:", err.Error())
		return err
	}

	if transaction.DeletedAt != nil {
		return fmt.Errorf("Transação %d já deletada.", id)
	}

	if transaction.Origin == transactionsconstants.IncomeReceiptOrigin {
		return fmt.Errorf("Transação %d originada de recebimento, excluir via recebimentos.", id)
	}

	if err := s.repository.Delete(ctx, id); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao deletar a transação:", err)
		return err
	}

	return nil
}
