package telegram_bot

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/controller/http/telegram_web_handler"
	"github.com/x0k/veterinary-clinic-backend/internal/controller/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
	"gopkg.in/telebot.v3"
)

const component_name = "telegram_bot"

type Config struct {
	CalendarWebAppOrigin     string
	CalendarInputHandlerPath string
	WebHandlerAddress        string
	Token                    string
	PollerTimeout            time.Duration
}

type Bot struct {
	httpService  *shared.HttpService
	cfg          *Config
	bot          *telebot.Bot
	clinic       *usecase.ClinicUseCase[shared.TelegramResponse]
	clinicDialog *usecase.ClinicDialogUseCase[shared.TelegramResponse]
}

func New(
	log *logger.Logger,
	clinic *usecase.ClinicUseCase[shared.TelegramResponse],
	clinicDialog *usecase.ClinicDialogUseCase[shared.TelegramResponse],
	fataler shared.Fataler,
	initDataParser telegram_web_handler.TelegramInitDataParser,
	cfg *Config,
) *Bot {
	mux := http.NewServeMux()
	telegram_web_handler.UseRouter(log, mux, clinicDialog, initDataParser, &telegram_web_handler.Config{
		CalendarWebAppOrigin:     cfg.CalendarWebAppOrigin,
		CalendarInputHandlerPath: cfg.CalendarInputHandlerPath,
	})
	return &Bot{
		cfg:          cfg,
		clinic:       clinic,
		clinicDialog: clinicDialog,
		httpService: shared.NewHttpService(
			component_name,
			&http.Server{
				Addr:    cfg.WebHandlerAddress,
				Handler: shared.Logging(log, mux),
			},
			fataler,
		),
	}
}

func (b *Bot) Name() string {
	return component_name
}

func (b *Bot) Start(ctx context.Context) error {
	const op = "infra.telegram_bot.Bot.Start"

	if err := b.httpService.Start(ctx); err != nil {
		return fmt.Errorf("%s starting http service: %w", op, err)
	}

	if bot, err := telebot.NewBot(telebot.Settings{
		Token: b.cfg.Token,
		Poller: &telebot.LongPoller{
			Timeout: b.cfg.PollerTimeout,
		},
	}); err != nil {
		return fmt.Errorf("%s starting telebot: %w", op, err)
	} else {
		b.bot = bot
	}
	telegram.UseRouter(ctx, b.bot, b.clinic, b.clinicDialog)
	go b.bot.Start()
	return nil
}

func (b *Bot) Stop(ctx context.Context) error {
	b.bot.Stop()
	return b.httpService.Stop(ctx)
}
