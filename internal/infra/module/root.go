package module

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
)

type Root struct {
	*Module
}

func NewRoot(log *logger.Logger) *Root {
	return &Root{
		Module: New(log, "root"),
	}
}

func (r *Root) awaiter(ctx context.Context) error {
	r.log.Info(ctx, "press CTRL-C to exit")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	select {
	case s := <-stop:
		r.log.Info(ctx, "received signal", slog.String("signal", s.String()))
		return nil
	case err := <-r.fatal:
		return err
	}
}

func (r *Root) Start(ctx context.Context) error {
	return r.start(ctx, r.awaiter)
}
