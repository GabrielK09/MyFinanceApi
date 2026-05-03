package transactionhelper

import (
	"context"
	loggerHelper "my_finance/internal/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func WithTransaction(
	ctx context.Context,
	db *pgxpool.Pool,
	fn func(tx pgx.Tx) error,
) error {
	tx, err := db.Begin(ctx)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao iniciar a transação:", err)
		return err
	}

	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao commitar a transação:", err)
		return err
	}

	return nil
}

func WithTransactionResult[T any](
	ctx context.Context,
	db *pgxpool.Pool,
	fn func(tx pgx.Tx) (T, error),
) (T, error) {
	var zero T

	tx, err := db.Begin(ctx)
	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao iniciar a transação:", err)
		return zero, err
	}

	defer tx.Rollback(ctx)

	result, err := fn(tx)

	if err != nil {
		return zero, err
	}

	if err := tx.Commit(ctx); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao commitar a transação:", err)
		return zero, err
	}

	return result, nil
}
