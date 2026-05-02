package transactionsrepository

import (
	"context"
	"errors"

	"my_finance/internal/apperrors"
	transactionsconstants "my_finance/internal/constants/transactions"
	loggerHelper "my_finance/internal/logger"
	transactions_model "my_finance/models/transactions"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionsRepository struct {
	db *pgxpool.Pool
}

func NewTransactionsRepository(db *pgxpool.Pool) *TransactionsRepository {
	return &TransactionsRepository{
		db: db,
	}
}

func (r *TransactionsRepository) GetAll(ctx context.Context) ([]transactions_model.TransactionsModel, error) {
	var transactions []transactions_model.TransactionsModel

	transactionsRows, err := r.db.Query(
		ctx,
		`
			SELECT
				id,
				origin,
				origin_id,
				category_id,
				type,
				description,
				amount,
				due_date,
				paid_at,
				status,
				notes,
				created_at,
				updated_at,
				deleted_at
			FROM
				transactions
			WHERE
				deleted_at IS NULL
			ORDER BY
				id;
		`,
	)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao executar o select:", err)
		return []transactions_model.TransactionsModel{}, err
	}

	defer transactionsRows.Close()

	for transactionsRows.Next() {
		var transaction transactions_model.TransactionsModel

		if err := transactionsRows.Scan(
			&transaction.Id,
			&transaction.Origin,
			&transaction.OriginId,
			&transaction.CategoryId,
			&transaction.Type,
			&transaction.Description,
			&transaction.Amount,
			&transaction.DueDate,
			&transaction.PaidAt,
			&transaction.Status,
			&transaction.Notes,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
			&transaction.DeletedAt,
		); err != nil {
			loggerHelper.ErrorLogger.Println("Erro ao ler os dados da consulta:", err)
			return []transactions_model.TransactionsModel{}, err
		}

		transactions = append(transactions, transaction)
	}

	if transactionsRows.Err() != nil {
		loggerHelper.ErrorLogger.Println("Erro ao ler os dados da consulta:", err)
		return []transactions_model.TransactionsModel{}, err
	}

	return transactions, nil
}

func (r *TransactionsRepository) Create(ctx context.Context, transaction transactions_model.TransactionsModel) (transactions_model.TransactionsModel, error) {
	var transactionData transactions_model.TransactionsModel

	tx, err := r.db.Begin(ctx)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao iniciar a transação:", err)
		return transactions_model.TransactionsModel{}, err
	}

	defer tx.Rollback(ctx)

	if err := tx.QueryRow(
		ctx,
		`
			INSERT INTO  transactions 
				(category_id, origin, origin_id, type, description, amount, due_date, paid_at, notes)
			VALUES 
				($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING 
				id, origin, origin_id, category_id, type, description, amount, due_date, paid_at, notes, status, created_at, updated_at
		`,
		transaction.CategoryId,
		transaction.Origin,
		transaction.OriginId,
		transaction.Type,
		transaction.Description,
		transaction.Amount,
		transaction.DueDate,
		transaction.PaidAt,
		transaction.Notes,
	).Scan(
		&transactionData.Id,
		&transactionData.Origin,
		&transactionData.OriginId,
		&transactionData.CategoryId,
		&transactionData.Type,
		&transactionData.Description,
		&transactionData.Amount,
		&transactionData.DueDate,
		&transactionData.PaidAt,
		&transactionData.Notes,
		&transactionData.Status,
		&transactionData.CreatedAt,
		&transactionData.UpdatedAt,
	); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao criar a transação:", err)
		return transactions_model.TransactionsModel{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao commitar:", err)
		return transactions_model.TransactionsModel{}, err
	}

	return transactionData, nil
}

func (r *TransactionsRepository) Update(ctx context.Context, transaction transactions_model.TransactionsModel) (transactions_model.TransactionsModel, error) {
	var transactionData transactions_model.TransactionsModel

	tx, err := r.db.Begin(ctx)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao iniciar a transação:", err)
		return transactions_model.TransactionsModel{}, err
	}

	defer tx.Rollback(ctx)

	err = tx.QueryRow(
		ctx,
		`
			UPDATE
				transactions 
			SET
				category_id = $2, 
				type = $3, 
				description = $4, 
				amount = $5, 
				due_date = $6, 
				paid_at = $7, 
				notes = $8
			WHERE
				id = $1
			RETURNING 
				id, origin, origin_id, category_id, type, description, amount, due_date, paid_at, notes, status, created_at, updated_at
		`,
		transaction.Id,
		transaction.CategoryId,
		transaction.Type,
		transaction.Description,
		transaction.Amount,
		transaction.DueDate,
		transaction.PaidAt,
		transaction.Notes,
	).Scan(
		&transactionData.Id,
		&transactionData.Origin,
		&transactionData.OriginId,
		&transactionData.CategoryId,
		&transactionData.Type,
		&transactionData.Description,
		&transactionData.Amount,
		&transactionData.DueDate,
		&transactionData.PaidAt,
		&transactionData.Notes,
		&transactionData.Status,
		&transactionData.CreatedAt,
		&transactionData.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			loggerHelper.ErrorLogger.Println("Trnsação não localizada:", err)
			return transactions_model.TransactionsModel{}, apperrors.ErrNotFound
		}

		loggerHelper.ErrorLogger.Println("Erro ao alterar a transação:", err)
		return transactions_model.TransactionsModel{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao commitar:", err)
		return transactions_model.TransactionsModel{}, err
	}

	return transactionData, nil
}

func (r *TransactionsRepository) FindById(ctx context.Context, id int) (transactions_model.TransactionsModel, error) {
	var transaction transactions_model.TransactionsModel

	err := r.db.QueryRow(
		ctx,
		`
			SELECT
				id,
				category_id,
				type,
				description,
				amount,
				due_date,
				paid_at,
				status,
				notes,
				created_at,
				updated_at,
				deleted_at
			FROM
				transactions
			WHERE
				id = $1
		`,
		id,
	).Scan(
		&transaction.Id,
		&transaction.CategoryId,
		&transaction.Type,
		&transaction.Description,
		&transaction.Amount,
		&transaction.DueDate,
		&transaction.PaidAt,
		&transaction.Status,
		&transaction.Notes,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&transaction.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			loggerHelper.ErrorLogger.Println("Trnsação não localizada:", err)
			return transactions_model.TransactionsModel{}, apperrors.ErrNotFound
		}

		loggerHelper.ErrorLogger.Println("Erro ao ler os dados da consulta:", err)
		return transactions_model.TransactionsModel{}, err
	}

	return transaction, nil
}

func (r *TransactionsRepository) Pay(ctx context.Context, id int) error {
	cmd, err := r.db.Exec(
		ctx,
		`
			UPDATE
				transactions
			SET
				paid_at = NOW(),
				status = $2,
				updated_at = NOW()
			WHERE
				id = $1 AND
				status = 'Pendente' AND
				paid_at IS NULL AND
				deleted_at IS NULL
		`,
		id,
		transactionsconstants.StatusPago,
	)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao pagar a transação:", err)
		return err
	}

	if cmd.RowsAffected() == 0 {
		loggerHelper.ErrorLogger.Println("Transação não localizada:", err)
		return apperrors.ErrNotFound
	}

	return nil
}

func (r *TransactionsRepository) Cancel(ctx context.Context, id int) error {
	cmd, err := r.db.Exec(
		ctx,
		`
			UPDATE
				transactions
			SET
				status = $2,
				updated_at = NOW()
			WHERE
				id = $1 AND
				paid_at IS NOT NULL AND
				deleted_at IS NULL
		`,
		id,
		transactionsconstants.StatusCancelado,
	)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao cancelar a transação:", err)
		return err
	}

	if cmd.RowsAffected() == 0 {
		loggerHelper.ErrorLogger.Println("Transação não localizada:", err)
		return apperrors.ErrNotFound
	}

	return nil
}

func (r *TransactionsRepository) Delete(ctx context.Context, id int) error {
	cmd, err := r.db.Exec(
		ctx,
		`
			UPDATE
				transactions
			SET
				deleted_at = NOW(),
				updated_at = NOW()
			WHERE
				id = $1 AND
				origin <> 'income_receipt' AND
				deleted_at IS NULL
		`,
		id,
	)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao cancelar a transação:", err)
		return err
	}

	if cmd.RowsAffected() == 0 {
		loggerHelper.ErrorLogger.Println("Transação não localizada:", err)
		return apperrors.ErrNotFound
	}

	return nil
}
