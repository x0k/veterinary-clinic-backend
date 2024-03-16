package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters/presenter"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters/repo"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters/repo/notion_repo"
	"github.com/x0k/veterinary-clinic-backend/internal/config"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/app_logger"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/boot"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/telegram_http_server"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"

	"github.com/x0k/veterinary-clinic-backend/internal/infra/profiler_http_server"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/telegram_bot"

	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
)

const calendarInputHandlerPath = "/calendar-input"

type CalendarRequestOptions struct {
	Url string `json:"url"`
}

func Run(cfg *config.Config) {
	ctx := context.Background()
	log := app_logger.MustNew(&cfg.Logger)
	if err := run(ctx, cfg, log); err != nil {
		log.Error(ctx, "failed to run", sl.Err(err))
	}
}

func run(ctx context.Context, cfg *config.Config, log *logger.Logger) error {

	log.Info(ctx, "starting application", slog.String("log_level", cfg.Logger.Level))

	b := boot.New(log)

	freePeriodsRepo := repo.NewHttpFreePeriods(
		log,
		cfg.ProductionCalendar.Url,
		&http.Client{},
	)

	calendarWebAppParams := url.Values{}
	calendarWebAppRequestOptions, err := json.Marshal(&CalendarRequestOptions{
		Url: fmt.Sprintf("%s%s", cfg.Telegram.WebHandlerUrl, calendarInputHandlerPath),
	})
	if err != nil {
		return err
	}
	calendarWebAppParams.Add("req", string(calendarWebAppRequestOptions))
	configuredCalendarWebAppUrl := fmt.Sprintf("%s?%s", cfg.Telegram.CalendarWebAppUrl, calendarWebAppParams.Encode())
	log.Debug(ctx, "configured calendar web app url", slog.String("url", configuredCalendarWebAppUrl))

	calendarWebAppUrl, err := url.Parse(cfg.Telegram.CalendarWebAppUrl)
	if err != nil {
		return err
	}
	calendarWebAppOrigin := fmt.Sprintf("%s://%s", calendarWebAppUrl.Scheme, calendarWebAppUrl.Host)

	notionClient := notionapi.NewClient(notionapi.Token(cfg.Notion.Token))
	clinicDialogUseCase := usecase.NewClinicDialogUseCase(
		log,
		presenter.NewTelegramDialog(&presenter.TelegramDialogConfig{
			CalendarWebAppUrl: configuredCalendarWebAppUrl,
		}),
		repo.NewStaticWorkBreaks(),
		notion_repo.NewBusyPeriods(
			log,
			notionClient,
			notionapi.DatabaseID(cfg.Notion.RecordsDatabaseId),
		),
		freePeriodsRepo,
	)

	b.Append(
		freePeriodsRepo,
		telegram_http_server.New(
			log,
			clinicDialogUseCase,
			&telegram_http_server.Config{
				Token:                    cfg.Telegram.Token,
				CalendarInputHandlerPath: calendarInputHandlerPath,
				CalendarWebAppOrigin:     calendarWebAppOrigin,
				Address:                  cfg.Telegram.WebHandlerAddress,
			},
		),
		telegram_bot.New(
			log,
			usecase.NewClinicUseCase(
				notion_repo.NewClinic(
					notionClient,
					notionapi.DatabaseID(cfg.Notion.ServicesDatabaseId),
					notionapi.DatabaseID(cfg.Notion.RecordsDatabaseId),
				),
				presenter.NewTelegramClinic(),
			),
			clinicDialogUseCase,
			&telegram_bot.Config{
				Token:         cfg.Telegram.Token,
				PollerTimeout: cfg.Telegram.PollerTimeout,
			},
		),
	)

	if cfg.Profiler.Enabled {
		b.Append(profiler_http_server.New(log, &cfg.Profiler))
	}

	b.Start(ctx)

	log.Info(ctx, "application stopped")
	return nil
}
