package boot

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

type Service interface {
	Name() string
	Start(ctx context.Context) error
	// Should not be called if the `Start` method returned an error
	Stop(ctx context.Context) error
}

type Boot struct {
	log                    *logger.Logger
	services               []Service
	fatal                  chan error
	fataled                atomic.Bool
	lastInitializedService int
}

func New(log *logger.Logger) *Boot {
	return &Boot{
		log:                    log.With(slog.String("component", "infra.boot.Boot")),
		fatal:                  make(chan error, 1),
		lastInitializedService: -1,
	}
}

func (b *Boot) Append(services ...Service) {
	b.services = append(b.services, services...)
}

func (b *Boot) TryAppend(ctx context.Context, services ...func() (Service, error)) {
	for _, f := range services {
		service, err := f()
		if err != nil {
			b.Fatal(ctx, err)
			break
		}
		b.services = append(b.services, service)
	}
}

func (b *Boot) Fatal(ctx context.Context, err error) {
	if b.fataled.Swap(true) {
		b.log.Error(ctx, "another fatal error", sl.Err(err))
		return
	}
	b.fatal <- err
	close(b.fatal)
}

func (b *Boot) Start(ctx context.Context) {
	if b.fataled.Load() {
		b.log.Error(ctx, "fataled before start", sl.Err(<-b.fatal))
		return
	}
	for i, service := range b.services {
		if err := service.Start(ctx); err != nil {
			b.Fatal(ctx, fmt.Errorf("failed to start %s: %w", service.Name(), err))
			break
		}
		b.log.Info(ctx, "started", slog.String("service", service.Name()))
		b.lastInitializedService = i
	}

	b.log.Info(ctx, "press CTRL-C to exit")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	select {
	case s := <-stop:
		b.log.Info(ctx, "received signal", slog.String("signal", s.String()))
	case err := <-b.fatal:
		b.log.Error(ctx, "fatal error", sl.Err(err))
	}

	b.log.Info(ctx, "shutting down")

	b.stop(ctx)
}

func (b *Boot) stop(ctx context.Context) {
	for i := b.lastInitializedService; i >= 0; i-- {
		service := b.services[i]
		if err := service.Stop(ctx); err != nil {
			b.log.Error(ctx, "failed to stop", slog.String("service", service.Name()), sl.Err(err))
		} else {
			b.log.Info(ctx, "stopped", slog.String("service", service.Name()))
		}
	}
}
