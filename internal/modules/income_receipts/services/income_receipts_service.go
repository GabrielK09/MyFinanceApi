package incomereceiptsservice

import (
	"context"
	"errors"
	"fmt"
	"my_finance/internal/apperrors"
	incomereceiptssconstants "my_finance/internal/constants/income_receipts"
	loggerHelper "my_finance/internal/logger"
	incomereceiptsmodel "my_finance/models/income_receipts"
	incomesourcesmodel "my_finance/models/income_sources"
)

type IncomeReceiptsRepository interface {
	GetAll(ctx context.Context) ([]incomereceiptsmodel.IncomeReceiptsModel, error)
	Create(ctx context.Context, incomeReceipts incomereceiptsmodel.IncomeReceiptsModel) (incomereceiptsmodel.IncomeReceiptsModel, error)
	Cancel(ctx context.Context, id int) error
	Delete(ctx context.Context, id int) error
	FindById(ctx context.Context, id int) (*incomereceiptsmodel.IncomeReceiptsModel, error)
}

type IncomeSourcesRepository interface {
	FindIncomeSourceById(ctx context.Context, id int) (*incomesourcesmodel.IncomeSourcesModel, error)
}

type IncomeReceiptsService struct {
	repository              IncomeReceiptsRepository
	incomeSourcesRepository IncomeSourcesRepository
}

func NewIncomeReceiptsService(repository IncomeReceiptsRepository, incomeSourcesRepository IncomeSourcesRepository) *IncomeReceiptsService {
	return &IncomeReceiptsService{
		repository:              repository,
		incomeSourcesRepository: incomeSourcesRepository,
	}
}

func (s *IncomeReceiptsService) GetAll(ctx context.Context) ([]incomereceiptsmodel.IncomeReceiptsModel, error) {
	incomeReceipts, err := s.repository.GetAll(ctx)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao retornar todos os recebimentos.")
		return []incomereceiptsmodel.IncomeReceiptsModel{}, err
	}

	return incomeReceipts, nil
}

func (s *IncomeReceiptsService) Create(ctx context.Context, incomeReceipts incomereceiptsmodel.IncomeReceiptsModel) (incomereceiptsmodel.IncomeReceiptsModel, error) {
	if err := incomeReceipts.Validate(); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao validar os dados para criar o recebimento:", err)
		return incomereceiptsmodel.IncomeReceiptsModel{}, err
	}

	_, err := s.incomeSourcesRepository.FindIncomeSourceById(ctx, incomeReceipts.IncomeSourceId)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return incomereceiptsmodel.IncomeReceiptsModel{}, fmt.Errorf("Fonte de renda %d não localizada.", incomeReceipts.IncomeSourceId)
		}

		loggerHelper.ErrorLogger.Println("Erro ao localizar a fonte de renda do recebimento:", err)
		return incomereceiptsmodel.IncomeReceiptsModel{}, err
	}

	incomeReceipt, err := s.repository.Create(ctx, incomeReceipts)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao criar o recebimento:", err)
		return incomereceiptsmodel.IncomeReceiptsModel{}, err
	}

	return incomeReceipt, nil
}

func (s *IncomeReceiptsService) Cancel(ctx context.Context, id int) error {
	incomeReceipt, err := s.repository.FindById(ctx, id)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			loggerHelper.ErrorLogger.Printf("Recebimento %d não localizada.", id)

			return apperrors.NewErrNotFound(fmt.Sprintf("Recebimento %d não localizada.", id))

		}

		loggerHelper.ErrorLogger.Println("Erro ao localizar o recebimento:", err)
		return err
	}

	if incomeReceipt.Status == incomereceiptssconstants.StatusCancelado {
		return fmt.Errorf("Recebimento %d já cancelado.", id)
	}

	if incomeReceipt.DeletedAt != nil {
		return fmt.Errorf("Recebimento %d deletado.", id)
	}

	if err := s.repository.Cancel(ctx, id); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao cancelar o recebimento:", err)
		return err
	}

	return nil
}

func (s *IncomeReceiptsService) Delete(ctx context.Context, id int) error {
	receipt, err := s.repository.FindById(ctx, id)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao localizar o recebimento para deletar:", err)
		return err
	}

	if receipt.DeletedAt != nil {
		loggerHelper.ErrorLogger.Println("Recebimento já deletado:", err)
		return fmt.Errorf("O recebimento n° %d já está deletado.", id)
	}

	return s.repository.Delete(ctx, id)
}

func (s *IncomeReceiptsService) FindById(ctx context.Context, id int) (*incomereceiptsmodel.IncomeReceiptsModel, error) {
	incomeReceipt, err := s.repository.FindById(ctx, id)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			loggerHelper.ErrorLogger.Printf("Recebimento %d não localizada.", id)

			return nil, apperrors.NewErrNotFound(fmt.Sprintf("Recebimento %d não localizada.", id))

		}

		loggerHelper.ErrorLogger.Println("Erro ao localizar o recebimento:", err)
		return nil, err
	}

	return incomeReceipt, nil
}
