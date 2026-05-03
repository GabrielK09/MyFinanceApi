package incomesourcesrepository

import (
	"context"
	"errors"
	"my_finance/internal/apperrors"
	constantsdbcode "my_finance/internal/constants/db"
	loggerHelper "my_finance/internal/logger"
	incomesourcesmodel "my_finance/models/income_sources"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IncomeSourcesRepository struct {
	db *pgxpool.Pool
}

func NewIncomeSourcesRepository(db *pgxpool.Pool) *IncomeSourcesRepository {
	return &IncomeSourcesRepository{
		db: db,
	}
}

func (r *IncomeSourcesRepository) GetAll(ctx context.Context) ([]incomesourcesmodel.IncomeSourcesModel, error) {
	incomeSourcesRows, err := r.db.Query(
		ctx,
		`
			SELECT
				id,
				name,
				description,
				active,
				created_at,
				updated_at
			FROM
				income_sources
			WHERE
				deleted_at IS NULL
		`,
	)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao executar o select:", err)
		return []incomesourcesmodel.IncomeSourcesModel{}, nil
	}

	defer incomeSourcesRows.Close()

	var incomeSources []incomesourcesmodel.IncomeSourcesModel

	for incomeSourcesRows.Next() {
		var incomeSource incomesourcesmodel.IncomeSourcesModel

		if err := incomeSourcesRows.Scan(
			&incomeSource.Id,
			&incomeSource.Name,
			&incomeSource.Description,
			&incomeSource.Active,
			&incomeSource.CreatedAt,
			&incomeSource.UpdatedAt,
		); err != nil {
			loggerHelper.ErrorLogger.Println("Erro ao ler os dados da consulta:", err)
			return []incomesourcesmodel.IncomeSourcesModel{}, err
		}

		incomeSources = append(incomeSources, incomeSource)
	}

	if err = incomeSourcesRows.Err(); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao interar os dados:", err)
		return incomeSources, err
	}

	return incomeSources, nil
}

func (r *IncomeSourcesRepository) Create(ctx context.Context, incomeSources incomesourcesmodel.IncomeSourcesModel) (incomesourcesmodel.IncomeSourcesModel, error) {
	var pgErr *pgconn.PgError
	var incomeSource incomesourcesmodel.IncomeSourcesModel

	err := r.db.QueryRow(
		ctx,
		`	
			INSERT INTO income_sources 
				(name, description)
			VALUES
				($1, $2)
			RETURNING
				id,
				name,
				description,
				created_at,
				updated_at,
				deleted_at
		`,
		incomeSources.Name,
		incomeSources.Description,
	).Scan(
		&incomeSource.Id,
		&incomeSource.Name,
		&incomeSource.Description,
		&incomeSource.CreatedAt,
		&incomeSource.UpdatedAt,
		&incomeSource.DeletedAt,
	)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao inserir a fonte de renda:", err)

		if errors.As(err, &pgErr) {

			if pgErr.Code == constantsdbcode.UniqueViolationCode {
				loggerHelper.ErrorLogger.Println(pgErr)
				loggerHelper.ErrorLogger.Println(apperrors.ErrUniqueConstraint.Error())

				return incomesourcesmodel.IncomeSourcesModel{}, apperrors.ErrUniqueConstraint
			}

			if pgErr.Code == constantsdbcode.ForeignKeyViolation {
				loggerHelper.ErrorLogger.Println(pgErr)
				loggerHelper.ErrorLogger.Println(apperrors.ErrForeignKeyViolation.Error())

				return incomesourcesmodel.IncomeSourcesModel{}, apperrors.ErrForeignKeyViolation
			}
		}

		return incomesourcesmodel.IncomeSourcesModel{}, err
	}

	return incomeSource, nil
}

func (r *IncomeSourcesRepository) Update(ctx context.Context, incomeSources incomesourcesmodel.IncomeSourcesModel) (incomesourcesmodel.IncomeSourcesModel, error) {
	var updatedIncomeSources incomesourcesmodel.IncomeSourcesModel

	tx, err := r.db.Begin(ctx)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao iniciar a transação:", err)
		return incomesourcesmodel.IncomeSourcesModel{}, err
	}

	defer tx.Rollback(ctx)

	if err := tx.QueryRow(
		ctx,
		`	
			UPDATE
				income_sources
			SET
				name = $2, 
				description = $3,
				updated_at = NOW()
			WHERE
				id = $1 AND
				deleted_at IS NULL
			RETURNING
				id,
				name,
				description,
				created_at,
				updated_at,
				deleted_at
		`,
		incomeSources.Id,
		incomeSources.Name,
		incomeSources.Description,
	).Scan(
		&updatedIncomeSources.Id,
		&updatedIncomeSources.Name,
		&updatedIncomeSources.Description,
		&updatedIncomeSources.CreatedAt,
		&updatedIncomeSources.UpdatedAt,
		&updatedIncomeSources.DeletedAt,
	); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao altera os dados da fonte de renda:", err)
		return incomesourcesmodel.IncomeSourcesModel{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao commitar os dados da fonte de renda:", err)
		return incomesourcesmodel.IncomeSourcesModel{}, err
	}

	return updatedIncomeSources, nil
}

func (r *IncomeSourcesRepository) Delete(ctx context.Context, id int) error {
	if _, err := r.db.Exec(
		ctx,
		`
			UPDATE
				income_sources
			SET	
				deleted_at = NOW(),
				updated_at = NOW(),
				active = false
			WHERE
				id = $1
		`,
		id,
	); err != nil {
		return err
	}

	return nil
}

func (r *IncomeSourcesRepository) Active(ctx context.Context, id int) error {
	if _, err := r.db.Exec(
		ctx,
		`
			UPDATE
				income_sources
			SET	
				deleted_at = null,
				updated_at = NOW(),
				active = true
			WHERE
				id = $1
		`,
		id,
	); err != nil {
		return err
	}

	return nil
}

func (r *IncomeSourcesRepository) FindById(ctx context.Context, id int) (*incomesourcesmodel.IncomeSourcesModel, error) {
	var incomeSource incomesourcesmodel.IncomeSourcesModel

	err := r.db.QueryRow(
		ctx,
		`
			SELECT
				id,
				name,
				description,
				active,
				created_at,
				updated_at,
				deleted_at
			FROM
				income_sources
			WHERE
				id = $1
		`,
		id,
	).Scan(
		&incomeSource.Id,
		&incomeSource.Name,
		&incomeSource.Description,
		&incomeSource.Active,
		&incomeSource.CreatedAt,
		&incomeSource.UpdatedAt,
		&incomeSource.DeletedAt,
	)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao localziar os dados da fonte de renda:", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}

		return nil, err
	}

	return &incomeSource, nil
}
