package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/x0k/veterinary-clinic-backend/internal/config"
	"github.com/x0k/veterinary-clinic-backend/internal/controller/profiler"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/app_logger"
)

func Run(cfg *config.Config) {
	log := app_logger.MustNew(&cfg.Logger)

	ctx := context.Background()

	log.Info(ctx, "starting application", slog.String("log_level", cfg.Logger.Level))

	prof := profiler.New(&cfg.Profiler, log)
	prof.Start(ctx)

	log.Info(ctx, "application started")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Info(ctx, "graceful shutdown")

	prof.Stop(ctx)

	log.Info(ctx, "application stopped")
}
