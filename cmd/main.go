package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/sync/errgroup"

	"github.com/madhurikadam/app-transcation/cmd/configuration"
	"github.com/madhurikadam/app-transcation/internal/database/postgres"
	httpGW "github.com/madhurikadam/app-transcation/internal/gateway/http"
	"github.com/madhurikadam/app-transcation/internal/service"
	postgresPkg "github.com/madhurikadam/app-transcation/pkg/database/postgres"
	httpPkg "github.com/madhurikadam/app-transcation/pkg/http"
	log "github.com/sirupsen/logrus"
)

const ServiceName = "app-transcation"

var cfg configuration.Config

func main() {

	// Parse the config from the environment, exiting on error
	if err := envconfig.Process("", &cfg); err != nil {
		log.Error("failed to parse config from environment", err)
		panic(err)
	}

	log.Info(fmt.Sprintf("starting %s service", ServiceName))
	ctx := context.Background()

	pgxPool, err := initPostgres(ctx)
	if err != nil {
		panic(err)
	}
	defer pgxPool.Close()

	repo := postgres.NewRepo(pgxPool, postgres.SetupPSQL())

	transcationSvc := service.New(&repo)

	errGroup, ctx := errgroup.WithContext(ctx)

	httpSvr, err := initHttpServer(&transcationSvc)
	if err != nil {
		panic(err)
	}

	errGroup.Go(func() error {
		<-ctx.Done()
		tCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		log.Info("attempting http server shutdown")

		return httpSvr.Shutdown(tCtx)()
	})

	errGroup.Go(func() error {
		err := httpSvr.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	if err := errGroup.Wait(); err != nil {
		panic(err)
	}

}

func initPostgres(ctx context.Context) (*pgxpool.Pool, error) {
	log.WithField("config", cfg).Info("connecting to postgres")

	pgxPool, err := postgresPkg.Wait(ctx, cfg.Config, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection to postgresPkg: %w", err)
	}

	log.Info("connected, running migrations...")
	if err = postgresPkg.Migrate(ctx, postgres.Migrations, "migrations", cfg.PostgresDSN()); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	return pgxPool, err
}

func initHttpServer(transcationSvc *service.TranscationService) (*httpPkg.Server, error) {
	gw := httpGW.NewGateway(transcationSvc)
	router := mux.NewRouter()

	router.HandleFunc("/accounts", gw.CreateAccount).Methods(http.MethodPost)
	router.HandleFunc("/accounts/{id:[-0-9a-zA-Z]+}", gw.GetAccount).Methods(http.MethodGet)

	router.HandleFunc("/transcations", gw.CreateTranscation).Methods(http.MethodPost)

	server := httpPkg.New(fmt.Sprintf(":%d", cfg.HTTPPort), router)

	return server, nil
}
