package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/http2"

	"golang.org/x/net/http2/h2c"
)

type Server struct {
	*http.Server
}

func New(addr string, h http.Handler) *Server {
	srv := &http.Server{
		Addr:         addr,
		Handler:      h2c.NewHandler(h, &http2.Server{}),
		ReadTimeout:  time.Duration(5 * time.Second),
		IdleTimeout:  time.Duration(30 * time.Second),
		WriteTimeout: time.Duration(5 * time.Second),
	}

	return &Server{srv}
}

func (s *Server) Start(ctx context.Context) func() error {
	return func() error {

		log.WithField("addr", s.Addr).Info("starting http server")
		err := s.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithField("addr", s.Addr).Error("failed to strat http server")
			return err
		}

		return nil
	}
}

func (s *Server) Shutdown(ctx context.Context) func() error {

	return func() error {
		<-ctx.Done()
		tCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Info("attempting http server shutdown")
		return s.Server.Shutdown(tCtx)
	}
}
