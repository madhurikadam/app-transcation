package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/madhurikadam/app-transcation/internal/domain"
)

type Repo struct {
	psql squirrel.StatementBuilderType
	pgx  *pgxpool.Pool
}

func NewRepo(pgx *pgxpool.Pool, psql squirrel.StatementBuilderType) Repo {
	return Repo{
		pgx:  pgx,
		psql: psql,
	}
}

func (r *Repo) CreateAccount(ctx context.Context, account domain.Account) error {

	stmt := r.psql.
		Insert(TableAccounts).
		Columns(
			ID,
			DocumentNumber,
			CreditLimit,
			WithdrewalLimit,
			CreatedAt,
			UpdatedAt,
		).
		Values(
			account.ID,
			account.DocumentNumber,
			account.CreaditLimit,
			account.WithdrawalLimit,
			account.CreatedAt,
			account.UpdatedAt,
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

func (r Repo) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	stmt := r.psql.
		Select(
			ID,
			DocumentNumber,
			CreditLimit,
			WithdrewalLimit,
			CreatedAt,
			UpdatedAt,
		).
		From(TableAccounts).
		Where(squirrel.Eq{ID: id})

	query, params, err := stmt.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var account domain.Account

	err = r.pgx.
		QueryRow(ctx, query, params...).
		Scan(
			&account.ID,
			&account.DocumentNumber,
			&account.CreaditLimit,
			&account.WithdrawalLimit,
			&account.CreatedAt,
			&account.UpdatedAt,
		)

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *Repo) updateDebitLimit(ctx context.Context, accountID string, amount float64, tx pgx.Tx) error {
	stmt := r.psql.
		Update(TableAccounts).
		Set(WithdrewalLimit,
			squirrel.Expr("withdrawal_limit + ?", amount),
		).Where(squirrel.Eq{ID: accountID})

	query, params, err := stmt.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if _, err := tx.Exec(ctx, query, params...); err != nil {
		return err
	}

	return nil
}

func (r *Repo) updateCreditLimit(ctx context.Context, accountID string, amount float64, tx pgx.Tx) error {
	stmt := r.psql.
		Update(TableAccounts).
		Set(CreditLimit,
			squirrel.Expr("credit_limit - ?", amount),
		).Where(squirrel.Eq{ID: accountID})

	query, params, err := stmt.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if _, err := tx.Exec(ctx, query, params...); err != nil {
		return err
	}

	return nil
}
