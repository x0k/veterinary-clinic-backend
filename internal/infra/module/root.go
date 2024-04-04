package module

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
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
		r.log.Error(ctx, "fatal error", sl.Err(err))
		return nil
	}
}

func (r *Root) Start(ctx context.Context) {
	_ = r.start(ctx, r.awaiter)
}
