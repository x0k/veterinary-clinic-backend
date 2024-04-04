package app

import (
	"context"
	"net/http"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters/controller"
	appointment_module "github.com/x0k/veterinary-clinic-backend/internal/appointment/module"
	"github.com/x0k/veterinary-clinic-backend/internal/infra"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/module"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"gopkg.in/telebot.v3"
)

func NewRoot(cfg *Config, log *logger.Logger) (*module.Root, error) {
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
