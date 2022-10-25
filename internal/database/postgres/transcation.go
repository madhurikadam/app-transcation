package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/madhurikadam/app-transcation/internal/domain"
)

func (r *Repo) CreateCreditTranscation(ctx context.Context, transcation domain.Transcation, dbTxList []domain.DebitTx) error {
	tx, err := r.pgx.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction")
	}

	for _, update := range dbTxList {
		err := r.updateBalance(ctx, update.ID, update.Amount, tx)
		if err != nil {
			txErr := tx.Rollback(ctx)
			if txErr != nil {
				return txErr
			}

			return err
		}
	}

	if err := r.createTranscation(ctx, transcation, tx); err != nil {
		return err
	}

	if err := r.updateCreditLimit(ctx, transcation.AccountID, transcation.Amount, tx); err != nil {
		txErr := tx.Rollback(ctx)
		if txErr != nil {
			return txErr
		}

		return err
	}

	return tx.Commit(ctx)
}

func (r *Repo) CreateDebitTranscation(ctx context.Context, transcation domain.Transcation) error {
	tx, err := r.pgx.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction")
	}

	if err := r.createTranscation(ctx, transcation, tx); err != nil {
		return err
	}

	if err := r.updateDebitLimit(ctx, transcation.AccountID, transcation.Amount, tx); err != nil {
		txErr := tx.Rollback(ctx)
		if txErr != nil {
			return txErr
		}

		return err
	}

	return tx.Commit(ctx)
}

func (r *Repo) createTranscation(ctx context.Context, transcation domain.Transcation, tx pgx.Tx) error {
	stmt := r.psql.
		Insert(TableTranscations).
		Columns(
			ID,
			AccountID,
			OperationTypeID,
			Amount,
			EventAt,
			Balance,
		).
		Values(
			transcation.ID,
			transcation.AccountID,
			transcation.OperationTypeID,
			transcation.Amount,
			transcation.EventAt,
			transcation.Balance,
		)

	query, params, err := stmt.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if _, err := tx.Exec(ctx, query, params...); err != nil {
		return err
	}

	return nil
}

func (r *Repo) ListDebitTx(ctx context.Context) ([]domain.Transcation, error) {
	return nil, nil
}
