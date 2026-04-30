package transactionsrepository

import (
	"context"
	"time"

	"my_finance/internal/constants"
	loggerHelper "my_finance/internal/logger"
	transactions_model "my_finance/models/transactions"

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

	return []transactions_model.TransactionsModel{}, nil
}

func (r *TransactionsRepository) Create(ctx context.Context, transaction transactions_model.TransactionsModel) (int, error) {
	var transactionId int

	if err := r.db.QueryRow(
		ctx,
		`
			INSERT INTO 
				transactions (category_id, type, description, amount, due_date, paid_at, notes)
			VALUES 
				($1, $2, $3, $4, $5, $6, $7)
			RETURNING 
				id
		`,
		transaction.CategoryId,
		transaction.Type,
		transaction.Description,
		transaction.Amount,
		transaction.DueDate,
		transaction.PaidAt,
		transaction.Notes,
	).Scan(&transactionId); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao iniciar a transação:", err)
		return 0, err
	}

	return transactionId, nil
}

func (r *TransactionsRepository) FindById(ctx context.Context, id int) (transactions_model.TransactionsModel, error) {
	var transaction transactions_model.TransactionsModel

	if err := r.db.QueryRow(
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
	); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao ler os dados da consulta:", err)
		return transactions_model.TransactionsModel{}, err
	}

	return transaction, nil
}

func (r *TransactionsRepository) Pay(ctx context.Context, id int, paidAt time.Time) error {
	if _, err := r.db.Exec(
		ctx,
		`
			UPDATE
				transactions
			SET
				paid_at = $2,
				status = $3,
				updated_at = NOW()
			WHERE
				id = $1 AND
				status = 'Pendente' AND
				paid_at IS NULL AND
				deleted_at IS NULL
		`,
		id,
		paidAt,
		constants.StatusPago,
	); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao pagar a transação:", err)
		return err
	}

	return nil
}
