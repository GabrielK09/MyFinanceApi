package incomereceiptsrepository

import (
	"context"
	"errors"
	"my_finance/internal/apperrors"
	loggerHelper "my_finance/internal/logger"
	incomereceiptsmodel "my_finance/models/income_receipts"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IncomeReceiptsRepository struct {
	db *pgxpool.Pool
}

func NewIncomeReceiptsRepository(db *pgxpool.Pool) *IncomeReceiptsRepository {
	return &IncomeReceiptsRepository{
		db: db,
	}
}

func (r *IncomeReceiptsRepository) GetAll(ctx context.Context) ([]incomereceiptsmodel.IncomeReceiptsModel, error) {
	var incomeReceipts []incomereceiptsmodel.IncomeReceiptsModel

	incomeReceiptsRows, err := r.db.Query(
		ctx,
		`
			SELECT
				id,
				income_source_id,
				description,
				amount,
				received_at,
				reference_month,
				reference_year,
				notes,
				created_at,
				updated_at
			FROM
				income_receipts
			WHERE
				deleted_at IS NULL
		`,
	)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao executar o select:", err)
		return []incomereceiptsmodel.IncomeReceiptsModel{}, err
	}

	defer incomeReceiptsRows.Close()

	for incomeReceiptsRows.Next() {
		var incomeReceipt incomereceiptsmodel.IncomeReceiptsModel

		if err := incomeReceiptsRows.Scan(
			&incomeReceipt.Id,
			&incomeReceipt.IncomeSourceId,
			&incomeReceipt.Description,
			&incomeReceipt.Amount,
			&incomeReceipt.ReceivedAt,
			&incomeReceipt.ReferenceMonth,
			&incomeReceipt.ReferenceYear,
			&incomeReceipt.Notes,
			&incomeReceipt.CreatedAt,
			&incomeReceipt.UpdatedAt,
		); err != nil {
			loggerHelper.ErrorLogger.Println("Erro ao ler os dados do recebimento:", err)
			return []incomereceiptsmodel.IncomeReceiptsModel{}, err
		}

		incomeReceipts = append(incomeReceipts, incomeReceipt)
	}

	return incomeReceipts, nil
}

func (r *IncomeReceiptsRepository) Create(ctx context.Context, incomeReceipts incomereceiptsmodel.IncomeReceiptsModel) (incomereceiptsmodel.IncomeReceiptsModel, error) {
	var incomeReceipt incomereceiptsmodel.IncomeReceiptsModel

	if err := r.db.QueryRow(
		ctx,
		`
			INSERT INTO income_receipts 
				(income_source_id, description, amount, received_at, reference_month, reference_year, notes)
			VALUES
				($1, $2, $3, $4, $5, $6, $7)
			RETURNING
				id,
				income_source_id,
				description, 
				amount, 
				received_at, 
				reference_month, 
				reference_year, 
				notes,
				created_at,
				updated_at,
				deleted_at
		`,
		incomeReceipts.IncomeSourceId,
		incomeReceipts.Description,
		incomeReceipts.Amount,
		incomeReceipts.ReceivedAt,
		incomeReceipts.ReferenceMonth,
		incomeReceipts.ReferenceYear,
		incomeReceipts.Notes,
	).Scan(
		&incomeReceipt.Id,
		&incomeReceipt.IncomeSourceId,
		&incomeReceipt.Description,
		&incomeReceipt.Amount,
		&incomeReceipt.ReceivedAt,
		&incomeReceipt.ReferenceMonth,
		&incomeReceipt.ReferenceYear,
		&incomeReceipt.Notes,
		&incomeReceipt.CreatedAt,
		&incomeReceipt.UpdatedAt,
		&incomeReceipt.DeletedAt,
	); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao criar o recebimento:", err)
		return incomereceiptsmodel.IncomeReceiptsModel{}, err

	}

	return incomeReceipt, nil
}

func (r *IncomeReceiptsRepository) Cancel(ctx context.Context, id int) error {

	return nil
}

func (r *IncomeReceiptsRepository) Delete(ctx context.Context, id int) error {
	cmd, err := r.db.Exec(
		ctx,
		`
			UPDATE
				income_receipts
			SET
				deleted_at = NOW(),
				updated_at = NOW()
			WHERE
				id = $1 AND
				deleted_at IS NULL
		`,
		id,
	)

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao deletar o recebimento:", err)
		return err
	}

	if cmd.RowsAffected() == 0 {
		loggerHelper.ErrorLogger.Println("Recebimento não localizado:", err)
		return apperrors.ErrNotFound
	}

	return nil
}

func (r *IncomeReceiptsRepository) FindById(ctx context.Context, id int) (*incomereceiptsmodel.IncomeReceiptsModel, error) {
	var incomeReceipt incomereceiptsmodel.IncomeReceiptsModel

	err := r.db.QueryRow(
		ctx,
		`
			SELECT
				id,
				income_source_id,
				description,
				amount,
				received_at,
				reference_month,
				reference_year,
				notes,
				created_at,
				updated_at
			FROM
				income_receipts
			WHERE
				id = $1 AND
				deleted_at IS NULL
		`,
		id,
	).Scan(
		&incomeReceipt.Id,
		&incomeReceipt.IncomeSourceId,
		&incomeReceipt.Description,
		&incomeReceipt.Amount,
		&incomeReceipt.ReceivedAt,
		&incomeReceipt.ReferenceMonth,
		&incomeReceipt.ReferenceYear,
		&incomeReceipt.Notes,
		&incomeReceipt.CreatedAt,
		&incomeReceipt.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}

		return nil, err
	}

	return &incomeReceipt, nil

}
