package app

import (
	"log/slog"

	"github.com/jomei/notionapi"
	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
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
		Token: cfg.Telegram.Token.String(),
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
	m.Append(telegram_adapters.NewService("telegram_bot", bot))

	notion := notionapi.NewClient(cfg.Notion.Token)

	// Modules

	profilerModule := profiler_module.New(&cfg.Profiler, log)
	m.Append(profilerModule)

	telegramInitDataParser := telegram_adapters.NewInitDataParser(
		cfg.Telegram.Token,
		cfg.Telegram.InitDataExpiry,
	)

	appointmentModule, err := appointment_module.New(
		&cfg.Appointment,
		log,
		bot,
		notion,
		telegramInitDataParser,
	)
	if err != nil {
		return nil, err
	}
	m.Append(appointmentModule)

	return m, nil
}
