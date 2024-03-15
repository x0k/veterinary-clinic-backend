package boot

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

type Service interface {
	Start(ctx context.Context) error
}

type Boot struct {
	wg       sync.WaitGroup
	log      *logger.Logger
	services []Service
	fatal    chan error
	fataled  atomic.Bool
}

func New(log *logger.Logger) *Boot {
	return &Boot{
		log:   log.With(slog.String("component", "infra.boot.Boot")),
		fatal: make(chan error, 1),
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
		b.log.Error(ctx, "fatal error", sl.Err(err))
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
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for _, service := range b.services {
		b.wg.Add(1)
		go func() {
			defer b.wg.Done()
			if err := service.Start(ctx); err != nil {
				b.Fatal(ctx, err)
			}
		}()
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
	b.fataled.Store(true)
	cancel()
	b.wg.Wait()
}
