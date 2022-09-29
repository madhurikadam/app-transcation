package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sethvargo/go-retry"

	log "github.com/sirupsen/logrus"
)

type ErrInvalidConfig struct {
	Err error
}

func (e *ErrInvalidConfig) Error() string {
	return fmt.Sprintf("failed to parse config: %s", e.Err.Error())
}

// Open attempts to connect to a Postgres database
func Open(ctx context.Context, cfg Config, afterConn func(context.Context, *pgx.Conn) error) (*pgxpool.Pool, error) {
	pgxCfg, err := pgxpool.ParseConfig(cfg.PostgresDSN())
	if err != nil {
		return nil, &ErrInvalidConfig{Err: err}
	}
	pgxCfg.AfterConnect = afterConn
	pgxCfg.MaxConns = cfg.MaxConnection

	pool, err := pgxpool.ConnectConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx conn pool: %w", err)
	}

	return pool, nil
}

// Wait is used to attempt establishing a connection to a Postgres database
func Wait(ctx context.Context, cfg Config, afterConn func(context.Context, *pgx.Conn) error) (*pgxpool.Pool, error) {

	var pool *pgxpool.Pool
	tryer := retry.NewFibonacci(1 * time.Second)

	err := retry.Do(ctx, retry.WithJitter(500*time.Millisecond, tryer), func(ctx context.Context) error {
		var err error
		pool, err = Open(ctx, cfg, afterConn)
		if err != nil {
			if isErrInvalidConfig(err) {
				log.Error("failed to open connection to postgres", err)
				return err
			}

			log.Warn("failed to open connection to postgres, retrying", err)
			return retry.RetryableError(err)
		}

		if err := pool.Ping(ctx); err != nil {
			log.Warn("failed to ping postgres, retrying", err)

			// This marks the error as retryable
			return retry.RetryableError(err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection to postgres: %w", err)
	}

	return pool, nil
}

func isErrInvalidConfig(err error) bool {
	var icErr *ErrInvalidConfig
	if ok := errors.As(err, &icErr); ok {
		return true
	}

	return false
}
