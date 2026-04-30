package categoriesservice

import (
	"context"
	"fmt"
	loggerHelper "my_finance/internal/logger"
	categories_model "my_finance/models/categories"
)

type CategoryRepository interface {
	GetAll(r context.Context) ([]categories_model.CategoryModel, error)
	Create(r context.Context, category categories_model.CategoryModel) (int, error)
	Update(r context.Context, category categories_model.CategoryModel) (categories_model.CategoryModel, error)
	FindById(r context.Context, id int) (categories_model.CategoryModel, error)
	Delete(r context.Context, id int) error
	Active(r context.Context, id int) error
}

type CategoryService struct {
	repository CategoryRepository
}

func NewCategoryService(repository CategoryRepository) *CategoryService {
	return &CategoryService{
		repository: repository,
	}
}

func (s *CategoryService) GetAll(ctx context.Context) ([]categories_model.CategoryModel, error) {
	categories, err := s.repository.GetAll(ctx)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao consultar todas as categorias:", err)
		return []categories_model.CategoryModel{}, err
	}

	return categories, nil
}

func (s *CategoryService) Create(ctx context.Context, category categories_model.CategoryModel) (int, error) {
	if err := category.Validate(); err != nil {
		return 0, err
	}

	categoryId, err := s.repository.Create(ctx, category)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao criar a categoria:", err)
		return 0, err
	}

	return categoryId, nil
}

func (s *CategoryService) Update(ctx context.Context, category categories_model.CategoryModel) (categories_model.CategoryModel, error) {
	category, err := s.repository.Update(ctx, category)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao alterar a categoria:", err)
		return categories_model.CategoryModel{}, err
	}

	return category, nil
}

func (s *CategoryService) FindById(ctx context.Context, id int) (categories_model.CategoryModel, error) {
	if id <= 0 {
		return categories_model.CategoryModel{}, fmt.Errorf("O ID da transação precisa ser maior que zero.")
	}

	category, err := s.repository.FindById(ctx, id)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao criar a categoria:", err)
		return categories_model.CategoryModel{}, err
	}

	return category, nil
}

func (s *CategoryService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("O ID da transação precisa ser maior que zero.")
	}

	if err := s.repository.Delete(ctx, id); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao deletar a categoria:", err)
		return err
	}

	return nil
}

func (s *CategoryService) Active(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("O ID da transação precisa ser maior que zero.")
	}

	if err := s.repository.Active(ctx, id); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao ativar a categoria:", err)
		return err
	}

	return nil
}
