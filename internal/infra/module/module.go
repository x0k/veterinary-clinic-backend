package module

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

type Service interface {
	Start(ctx context.Context) error
}

type Module struct {
	log      *logger.Logger
	wg       sync.WaitGroup
	services []Service
	fatal    chan error
	stopped  atomic.Bool
}

func New(log *logger.Logger, name string) *Module {
	return &Module{
		log: log.With(
			slog.String("component", "boot.Module"),
			slog.String("module_name", name),
		),
		fatal: make(chan error, 1),
	}
}

func (m *Module) Append(services ...Service) {
	m.services = append(m.services, services...)
}

func (m *Module) Fatal(ctx context.Context, err error) {
	if m.stopped.Swap(true) {
		m.log.Error(ctx, "fatal error", sl.Err(err))
		return
	}
	m.fatal <- err
	close(m.fatal)
}

func (m *Module) awaiter(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	case err := <-m.fatal:
		return err
	}
}

func (m *Module) start(ctx context.Context, awaiter func(context.Context) error) error {
	if m.stopped.Load() {
		return <-m.fatal
	}

	m.log.Info(ctx, "starting")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, service := range m.services {
		m.wg.Add(1)
		go func() {
			defer m.wg.Done()
			if err := service.Start(ctx); err != nil {
				m.Fatal(ctx, err)
			}
		}()
	}

	err := awaiter(ctx)

	m.log.Info(ctx, "shutting down")

	m.stopped.Store(true)
	cancel()

	m.wg.Wait()

	return err
}

func (m *Module) Start(ctx context.Context) error {
	return m.start(ctx, m.awaiter)
}
