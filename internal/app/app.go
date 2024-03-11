package app

import (
	"context"
	"log/slog"

	"github.com/x0k/veterinary-clinic-backend/internal/config"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/app_logger"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/boot"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/profiler_http_server"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/telegram_bot"
)

type Service interface {
	Start(ctx context.Context)
	Stop(ctx context.Context)
}

func Run(cfg *config.Config) {
	log := app_logger.MustNew(&cfg.Logger)

	ctx := context.Background()

	log.Info(ctx, "starting application", slog.String("log_level", cfg.Logger.Level))

	b := boot.New(log)

	b.Append(telegram_bot.New(cfg, log))

	if cfg.Profiler.Enabled {
		b.Append(profiler_http_server.New(&cfg.Profiler, log, b))
	}

	b.Start(ctx)

	log.Info(ctx, "application stopped")
}
