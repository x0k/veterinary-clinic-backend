package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/config"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/app_logger"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/boot"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/notion_clinic_repo"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/profiler_http_server"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/telegram_bot"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/telegram_clinic_presenter"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/telegram_dialog_presenter"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/telegram_init_data_parser"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
)

const calendarInputHandlerPath = "/calendar-input"

type CalendarRequestOptions struct {
	Url string `json:"url"`
}

func Run(cfg *config.Config) {
	log := app_logger.MustNew(&cfg.Logger)

	ctx := context.Background()

	log.Info(ctx, "starting application", slog.String("log_level", cfg.Logger.Level))

	b := boot.New(log)

	b.TryAppend(ctx, func() (boot.Service, error) {
		calendarWebAppParams := url.Values{}
		calendarWebAppRequestOptions, err := json.Marshal(&CalendarRequestOptions{
			Url: fmt.Sprintf("%s%s", cfg.Telegram.WebHandlerUrl, calendarInputHandlerPath),
		})
		if err != nil {
			return nil, err
		}
		calendarWebAppParams.Add("req", string(calendarWebAppRequestOptions))
		configuredCalendarWebAppUrl := fmt.Sprintf("%s?%s", cfg.Telegram.CalendarWebAppUrl, calendarWebAppParams.Encode())
		log.Debug(ctx, "configured calendar web app url", slog.String("url", configuredCalendarWebAppUrl))

		calendarWebAppUrl, err := url.Parse(cfg.Telegram.CalendarWebAppUrl)
		if err != nil {
			return nil, err
		}
		calendarWebAppOrigin := fmt.Sprintf("%s://%s", calendarWebAppUrl.Scheme, calendarWebAppUrl.Host)

		return telegram_bot.New(
			log,
			usecase.NewClinicUseCase(
				notion_clinic_repo.New(
					notionapi.NewClient(notionapi.Token(cfg.Notion.Token)),
					notionapi.DatabaseID(cfg.Notion.ServicesDatabaseId),
					notionapi.DatabaseID(cfg.Notion.RecordsDatabaseId),
				),
				telegram_clinic_presenter.New(),
			),
			usecase.NewClinicDialogUseCase(
				log,
				telegram_dialog_presenter.New(&telegram_dialog_presenter.Config{
					CalendarWebAppUrl: configuredCalendarWebAppUrl,
				}),
			),
			b,
			telegram_init_data_parser.New(cfg.Telegram.Token, 24*time.Hour),
			&telegram_bot.Config{
				CalendarInputHandlerPath: calendarInputHandlerPath,
				CalendarWebAppOrigin:     calendarWebAppOrigin,
				WebHandlerAddress:        cfg.Telegram.WebHandlerAddress,
				Token:                    cfg.Telegram.Token,
				PollerTimeout:            cfg.Telegram.PollerTimeout,
			},
		), nil
	})

	if cfg.Profiler.Enabled {
		b.Append(profiler_http_server.New(&cfg.Profiler, b))
	}

	b.Start(ctx)

	log.Info(ctx, "application stopped")
}
