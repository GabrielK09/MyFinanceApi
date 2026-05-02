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

func (i *IncomeSourcesService) GetAll(ctx context.Context) ([]incomesourcesmodel.IncomeSourcesModel, error) {
	incomeSources, err := i.repository.GetAll(ctx)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao consultar todas as fontes de renda:", err)
		return []incomesourcesmodel.IncomeSourcesModel{}, err
	}

	return incomeSources, nil
}

func (i *IncomeSourcesService) Create(ctx context.Context, incomeSource incomesourcesmodel.IncomeSourcesModel) (incomesourcesmodel.IncomeSourcesModel, error) {
	if err := incomeSource.Validate(); err != nil {
		return incomesourcesmodel.IncomeSourcesModel{}, err
	}

	incomeSources, err := i.repository.Create(ctx, incomeSource)

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

func (i *IncomeSourcesService) Update(ctx context.Context, incomeSource incomesourcesmodel.IncomeSourcesModel) (incomesourcesmodel.IncomeSourcesModel, error) {
	if err := incomeSource.Validate(); err != nil {
		return incomesourcesmodel.IncomeSourcesModel{}, err
	}

	updated, err := i.repository.Update(ctx, incomeSource)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao alterar a fonte de renda:", err)
		return incomesourcesmodel.IncomeSourcesModel{}, err
	}

	return updated, nil
}

func (i *IncomeSourcesService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("O ID da fonte de renda precisa ser maior que zero.")
	}

	if err := i.repository.Delete(ctx, id); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao deletar a fonte de renda:", err)
		return err
	}

	return nil
}

func (i *IncomeSourcesService) Active(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("O ID da fonte de renda precisa ser maior que zero.")
	}

	if err := i.repository.Active(ctx, id); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao ativar a fonte de renda:", err)
		return err
	}

	return nil
}
