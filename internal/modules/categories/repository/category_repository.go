package categoriesrepository

import (
	"context"
	"errors"
	"my_finance/internal/apperrors"
	loggerHelper "my_finance/internal/logger"
	categories_model "my_finance/models/categories"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (c *CategoryRepository) GetAll(ctx context.Context) ([]categories_model.CategoryModel, error) {
	var categories []categories_model.CategoryModel

	categoriesRows, err := c.db.Query(
		ctx,
		`
			SELECT 
				id,
				name,
				type,
				monthly_limit,
				created_at,
				updated_at,
				deleted_at
			FROM
				categories
		`,
	)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao executar o select:", err)
		return []categories_model.CategoryModel{}, err
	}

	defer categoriesRows.Close()

	for categoriesRows.Next() {
		var category categories_model.CategoryModel

		if err := categoriesRows.Scan(
			&category.Id,
			&category.Name,
			&category.Type,
			&category.MonthlyLimit,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.DeletedAt,
		); err != nil {
			loggerHelper.ErrorLogger.Println("Erro ao ler os dados do select:", err)
			return []categories_model.CategoryModel{}, err
		}

		categories = append(categories, category)

	}

	if err := categoriesRows.Err(); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao iterar os dados do select:", err)
		return []categories_model.CategoryModel{}, err
	}

	return categories, nil
}

func (c *CategoryRepository) Create(ctx context.Context, category categories_model.CategoryModel) (int, error) {
	var categoryId int

	if err := c.db.QueryRow(
		ctx,
		`
			INSERT INTO categories 
				(name, type, monthly_limit)
			VALUES
				($1, $2, $3)
			RETURNING 	
				id
		`,
		category.Name,
		category.Type,
		category.MonthlyLimit,
	).Scan(&categoryId); err != nil {

		return 0, err
	}

	return categoryId, nil
}

func (c *CategoryRepository) FindById(ctx context.Context, id int) (categories_model.CategoryModel, error) {
	var category categories_model.CategoryModel

	err := c.db.QueryRow(
		ctx,
		`
			SELECT 
				id,
				name,
				type,
				monthly_limit,
				created_at,
				updated_at,
				deleted_at
			FROM
				categories
			WHERE
				id = $1
		`,
		id,
	).Scan(
		&category.Id,
		&category.Name,
		&category.Type,
		&category.MonthlyLimit,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return categories_model.CategoryModel{}, apperrors.ErrNotFound
		}

		loggerHelper.ErrorLogger.Println("Erro ao ler os dados da consulta:", err)
		return categories_model.CategoryModel{}, err
	}

	return category, nil
}

func (c *CategoryRepository) Update(ctx context.Context, category categories_model.CategoryModel) (categories_model.CategoryModel, error) {
	var updatedCategory categories_model.CategoryModel

	tx, err := c.db.Begin(ctx)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao iniciar a transação:", err)
		return categories_model.CategoryModel{}, err
	}

	defer tx.Rollback(ctx)

	if err := tx.QueryRow(
		ctx,
		`
			UPDATE
				categories
			SET
				name = $2,
				type = $3,
				monthly_limit = $4,
				updated_at = NOW()
			WHERE
				id = $1
			RETURNING
				id,
				name,
				type,
				monthly_limit,
				created_at,
				updated_at,
				deleted_at
		`,
		category.Id,
		category.Name,
		category.Type,
		category.MonthlyLimit,
	).Scan(
		&updatedCategory.Id,
		&updatedCategory.Name,
		&updatedCategory.Type,
		&updatedCategory.MonthlyLimit,
		&updatedCategory.CreatedAt,
		&updatedCategory.UpdatedAt,
		&updatedCategory.DeletedAt,
	); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao alterar os dados da categoria:", err)
		return categories_model.CategoryModel{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao commitar as atualizações:", err)
		return categories_model.CategoryModel{}, err
	}

	return updatedCategory, nil
}

func (c *CategoryRepository) Delete(ctx context.Context, id int) error {
	deletedAt := time.Now()

	tx, err := c.db.Begin(ctx)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao iniciar a transação:", err)
		return err
	}

	defer tx.Rollback(ctx)

	if _, err := tx.Exec(
		ctx,
		`
			UPDATE
				categories
			SET
				deleted_at = $1
			WHERE
				id = $2	
		`,
		deletedAt,
		id,
	); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao deletar os a categoria:", err)
		return err
	}

	if _, err := tx.Exec(
		ctx,
		`
			UPDATE
				transactions
			SET
				deleted_at = (
					SELECT	
						deleted_at
					FROM
						categories
					WHERE
						id = $1
				)
			WHERE
				category_id = $1
		`,
		id,
	); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao deletar as transações associadas a categoria:", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao commitar as atualizações:", err)
		return err
	}

	return nil
}

func (c *CategoryRepository) Active(ctx context.Context, id int) error {
	deletedAt := time.Now()

	tx, err := c.db.Begin(ctx)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao iniciar a transação:", err)
		return err
	}

	defer tx.Rollback(ctx)

	if _, err := tx.Exec(
		ctx,
		`
			UPDATE
				categories
			SET
				updated_at = $1,
				deleted_at = NULL
			WHERE
				id = $2	
		`,
		deletedAt,
		id,
	); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao ativar a categoria:", err)
		return err
	}

	if _, err := tx.Exec(
		ctx,
		`
			UPDATE
				transactions
			SET
				deleted_at = NULL,
				updated_at = (
					SELECT	
						updated_at
					FROM
						categories
					WHERE
						id = $1
				)
			WHERE
				category_id = $1
		`,
		id,
	); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao ativar as transações associadas a categoria:", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao commitar as atualizações:", err)
		return err
	}

	return nil
}
