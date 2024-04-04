package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters/controller"
	appointment_module "github.com/x0k/veterinary-clinic-backend/internal/appointment/module"
	"github.com/x0k/veterinary-clinic-backend/internal/infra"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/app_logger"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/module"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"gopkg.in/telebot.v3"
)

func Run(cfg *Config) {
	ctx := context.Background()
	log := app_logger.MustNew(&cfg.Logger)
	log.Info(ctx, "starting application", slog.String("log_level", cfg.Logger.Level))
	root, err := newRoot(cfg, log)
	if err != nil {
		log.Error(ctx, "failed to run", sl.Err(err))
		return
	}
	if err := root.Start(ctx); err != nil {
		log.Error(ctx, "fatal error", sl.Err(err))
	}
	log.Info(ctx, "application stopped")
}

func newRoot(cfg *Config, log *logger.Logger) (*module.Root, error) {
	m := module.NewRoot(log)

	bot, err := telebot.NewBot(telebot.Settings{
		Token: string(cfg.Telegram.Token),
		Poller: &telebot.LongPoller{
			Timeout: cfg.Telegram.PollerTimeout,
		},
	})
	if err != nil {
		return nil, err
	}

	notion := notionapi.NewClient(cfg.Notion.Token)

	appointmentModule, err := appointment_module.New(
		&cfg.Appointment,
		log,
		bot,
		notion,
	)
	if err != nil {
		return nil, err
	}

	m.Append(
		appointmentModule,
		infra.Starter(func(ctx context.Context) error {
			context.AfterFunc(ctx, func() {
				bot.Stop()
			})
			bot.Start()
			return nil
		}),
	)

	if cfg.Profiler.Enabled {
		m.Append(infra.NewHttpService(
			log,
			&http.Server{
				Addr:    cfg.Profiler.Address,
				Handler: controller.UseHttpProfilerRouter(http.NewServeMux()),
			},
		))
	}

	return m, nil
}
