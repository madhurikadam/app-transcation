package postgres

import (
	"context"
	"fmt"

	"github.com/madhurikadam/app-transcation/internal/domain"
)

func (r *Repo) CreateTranscation(ctx context.Context, transcation domain.Transcation) error {
	stmt := r.psql.
		Insert(TableTranscations).
		Columns(
			ID,
			AccountID,
			OperationTypeID,
			Amount,
			EventAt,
		).
		Values(
			transcation.ID,
			transcation.AccountID,
			transcation.OperationTypeID,
			transcation.Amount,
			transcation.EventAt,
		)

	query, params, err := stmt.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if _, err := r.pgx.Exec(ctx, query, params...); err != nil {
		return err
	}

	return nil
}
