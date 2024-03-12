package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/config"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/app_logger"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/boot"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/memory_dialog_repo"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/notion_clinic_repo"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/profiler_http_server"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/telegram_bot"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/telegram_clinic_presenter"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/telegram_dialog_presenter"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/telegram_init_data_parser"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
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

	telegramDialogPresenter, err := telegram_dialog_presenter.New(&telegram_dialog_presenter.TelegramDialogPresenterConfig{
		CalendarHandlerUrl: cfg.Telegram.WebHandlerUrl,
	})
	if err != nil {
		b.Fatal(ctx, err)
	}

	b.Append(
		telegram_bot.New(
			log,
			&cfg.Telegram,
			usecase.NewClinicUseCase(
				notion_clinic_repo.New(
					notionapi.NewClient(notionapi.Token(cfg.Notion.Token)),
					notionapi.DatabaseID(cfg.Notion.ServicesDatabaseId),
					notionapi.DatabaseID(cfg.Notion.RecordsDatabaseId),
				),
				telegram_clinic_presenter.New(),
			),
			usecase.NewClinicDialogUseCase(
				memory_dialog_repo.New(),
				telegramDialogPresenter,
			),
			b,
			telegram_init_data_parser.New(cfg.Telegram.Token, 24*time.Hour),
		),
	)

	if cfg.Profiler.Enabled {
		b.Append(profiler_http_server.New(&cfg.Profiler, b))
	}

	b.Start(ctx)

	log.Info(ctx, "application stopped")
}
