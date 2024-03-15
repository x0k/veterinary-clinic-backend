package infra

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

type HttpService struct {
	log *logger.Logger
	srv *http.Server
}

func NewHttpService(log *logger.Logger, srv *http.Server) *HttpService {
	return &HttpService{
		log: log.With(slog.String("component", "infra.http_service.HttpService")),
		srv: srv,
	}
}

func (s *HttpService) Start(ctx context.Context) error {
	const op = "infra.HttpService.Start"
	context.AfterFunc(ctx, func() {
		if err := s.srv.Shutdown(ctx); err != nil {
			s.log.Error(ctx, "failed to shutdown http server", sl.Err(err))
		}
	})
	return s.srv.ListenAndServe()
}
