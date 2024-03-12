package shared

import (
	"context"
	"fmt"
	"net/http"
)

type HttpService struct {
	name    string
	srv     *http.Server
	fataler Fataler
}

func NewHttpService(name string, srv *http.Server, fataler Fataler) *HttpService {
	return &HttpService{
		name:    name,
		srv:     srv,
		fataler: fataler,
	}
}

func (s *HttpService) Start(ctx context.Context) error {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.fataler.Fatal(ctx, fmt.Errorf("%s failed to start: %w", s.name, err))
		}
	}()
	return nil
}

func (s *HttpService) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
