package http_server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

type HttpServer struct {
	srv *http.Server
	log *logger.Logger
}

func New(srv *http.Server, log *logger.Logger) *HttpServer {
	return &HttpServer{
		srv: srv,
		log: log.With(slog.String("component", "infra.http_server")),
	}
}

func (s *HttpServer) Start(ctx context.Context) {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error(ctx, "failed to start", sl.Err(err))
		}
	}()
}

func (s *HttpServer) Stop(ctx context.Context) {
	if err := s.srv.Shutdown(ctx); err != nil {
		s.log.Error(ctx, "failed to stop", sl.Err(err))
	}
}
