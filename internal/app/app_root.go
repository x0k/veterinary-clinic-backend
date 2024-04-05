package app

import (
	"log/slog"

	"github.com/jomei/notionapi"
	adapters_telegram "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	appointment_module "github.com/x0k/veterinary-clinic-backend/internal/appointment/module"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
	profiler_module "github.com/x0k/veterinary-clinic-backend/internal/profiler"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
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
	bot.Use(
		middleware.Logger(slog.NewLogLogger(log.Logger.Handler(), slog.LevelDebug)),
		middleware.AutoRespond(),
		middleware.Recover(),
	)
	m.Append(adapters_telegram.NewService("telegram_bot", bot))

	notion := notionapi.NewClient(cfg.Notion.Token)

	// Modules

	profilerModule := profiler_module.New(&cfg.Profiler, log)
	m.Append(profilerModule)

	appointmentModule, err := appointment_module.New(
		&cfg.Appointment,
		log,
		bot,
		notion,
	)
	if err != nil {
		return nil, err
	}
	m.Append(appointmentModule)

	return m, nil
}
