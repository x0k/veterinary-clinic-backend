package app

import (
	"github.com/jomei/notionapi"
	appointment_module "github.com/x0k/veterinary-clinic-backend/internal/appointment/module"
	adapters_telegram "github.com/x0k/veterinary-clinic-backend/internal/infra/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
	"github.com/x0k/veterinary-clinic-backend/internal/profiler"
	"gopkg.in/telebot.v3"
)

func NewRoot(cfg *Config, log *logger.Logger) (*module.Root, error) {
	m := module.NewRoot(log.Logger)

	// Infrastructure

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

	// Modules

	profilerModule := profiler.New(&cfg.Profiler, log)
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
		profilerModule,
		appointmentModule,
		adapters_telegram.NewService(bot, log),
	)

	return m, nil
}
