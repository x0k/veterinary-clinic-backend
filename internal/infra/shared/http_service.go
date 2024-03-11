package shared

import (
	"context"
	"net/http"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

type HttpService struct {
	log     *logger.Logger
	srv     *http.Server
	fataler Fataler
}

func NewHttpService(srv *http.Server, log *logger.Logger, fataler Fataler) *HttpService {
	return &HttpService{
		log:     log,
		srv:     srv,
		fataler: fataler,
	}
}

func (s *HttpService) Start(ctx context.Context) error {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error(ctx, "failed to start", sl.Err(err))
			s.fataler.Fatal(ctx, err)
		}
	}()
	return nil
}

func (s *HttpService) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
