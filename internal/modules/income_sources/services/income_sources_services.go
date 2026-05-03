package incomesourcesservices

import (
	"context"
	"errors"
	"fmt"
	"my_finance/internal/apperrors"
	loggerHelper "my_finance/internal/logger"
	incomesourcesmodel "my_finance/models/income_sources"
)

type IncomeSourcesRepository interface {
	GetAll(ctx context.Context) ([]incomesourcesmodel.IncomeSourcesModel, error)
	FindById(ctx context.Context, id int) (*incomesourcesmodel.IncomeSourcesModel, error)
	Create(ctx context.Context, incomeSource incomesourcesmodel.IncomeSourcesModel) (incomesourcesmodel.IncomeSourcesModel, error)
	Update(ctx context.Context, incomeSource incomesourcesmodel.IncomeSourcesModel) (incomesourcesmodel.IncomeSourcesModel, error)
	Delete(ctx context.Context, id int) error
	Active(ctx context.Context, id int) error
}

type IncomeSourcesService struct {
	repository IncomeSourcesRepository
}

func NewIncomeSourcesService(repository IncomeSourcesRepository) *IncomeSourcesService {
	return &IncomeSourcesService{
		repository: repository,
	}
}

func (s *IncomeSourcesService) GetAll(ctx context.Context) ([]incomesourcesmodel.IncomeSourcesModel, error) {
	incomeSources, err := s.repository.GetAll(ctx)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao consultar todas as fontes de renda:", err)
		return []incomesourcesmodel.IncomeSourcesModel{}, err
	}

	return incomeSources, nil
}

func (s *IncomeSourcesService) FindById(ctx context.Context, id int) (*incomesourcesmodel.IncomeSourcesModel, error) {
	incomeSource, err := s.repository.FindById(ctx, id)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			loggerHelper.ErrorLogger.Printf("Fonte de renda %d não localizada.", id)
			return nil, apperrors.NewErrNotFound(fmt.Sprintf("Fonte de renda %d não localizada.", id))
		}

		loggerHelper.ErrorLogger.Println("Erro ao localizar a fonte de renda:", err)
		return nil, err
	}

	return incomeSource, nil
}

func (s *IncomeSourcesService) Create(ctx context.Context, incomeSource incomesourcesmodel.IncomeSourcesModel) (incomesourcesmodel.IncomeSourcesModel, error) {
	if err := incomeSource.Validate(); err != nil {
		return incomesourcesmodel.IncomeSourcesModel{}, err
	}

	incomeSources, err := s.repository.Create(ctx, incomeSource)

	if err != nil {
		if errors.Is(err, apperrors.ErrUniqueConstraint) {
			loggerHelper.ErrorLogger.Println("ErrUniqueConstraint:", err)
			return incomesourcesmodel.IncomeSourcesModel{}, err
		}

		loggerHelper.ErrorLogger.Println("Erro ao criar a fonte de renda:", err)
		return incomesourcesmodel.IncomeSourcesModel{}, err
	}

	return incomeSources, nil
}

func (s *IncomeSourcesService) Update(ctx context.Context, incomeSource incomesourcesmodel.IncomeSourcesModel) (incomesourcesmodel.IncomeSourcesModel, error) {
	if err := incomeSource.Validate(); err != nil {
		return incomesourcesmodel.IncomeSourcesModel{}, err
	}

	incomeSourceData, err := s.FindById(ctx, incomeSource.Id)

	if incomeSourceData.DeletedAt != nil {
		return incomesourcesmodel.IncomeSourcesModel{}, fmt.Errorf("Fonte de renda %d deletada:", incomeSourceData.Id)
	}

	updated, err := s.repository.Update(ctx, incomeSource)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao alterar a fonte de renda:", err)
		return incomesourcesmodel.IncomeSourcesModel{}, err
	}

	return updated, nil
}

func (s *IncomeSourcesService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("O ID da fonte de renda precisa ser maior que zero.")
	}

	incomeSourceData, err := s.FindById(ctx, id)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.NewErrNotFound(fmt.Sprintf("Fonte de renda %d não localizada.", id))
		}

		return err
	}

	if incomeSourceData != nil && incomeSourceData.DeletedAt != nil {
		return fmt.Errorf("Fonte de renda %d já deletada", id)
	}

	if err := s.repository.Delete(ctx, id); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao deletar a fonte de renda:", err)
		return err
	}

	return nil
}

func (s *IncomeSourcesService) Active(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("O ID da fonte de renda precisa ser maior que zero.")
	}

	incomeSourceData, err := s.FindById(ctx, id)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.NewErrNotFound(fmt.Sprintf("Fonte de renda %d não localizada.", id))
		}

		return err
	}

	if incomeSourceData != nil && incomeSourceData.DeletedAt == nil {
		return fmt.Errorf("Fonte de renda %d já ativada", id)
	}

	if err := s.repository.Active(ctx, id); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao ativar a fonte de renda:", err)
		return err
	}

	return nil
}
