package telegram_bot

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
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
	log          *logger.Logger
	wg           sync.WaitGroup
	httpService  *adapters.HttpService
	cfg          *Config
	bot          *telebot.Bot
	clinic       *usecase.ClinicUseCase[adapters.TelegramResponse]
	clinicDialog *usecase.ClinicDialogUseCase[adapters.TelegramResponse]
}

func New(
	log *logger.Logger,
	clinic *usecase.ClinicUseCase[adapters.TelegramResponse],
	clinicDialog *usecase.ClinicDialogUseCase[adapters.TelegramResponse],
	fataler adapters.Fataler,
	initDataParser telegram_web_handler.TelegramInitDataParser,
	cfg *Config,
) *Bot {
	mux := http.NewServeMux()
	telegram_web_handler.UseRouter(log, mux, clinicDialog, initDataParser, &telegram_web_handler.Config{
		CalendarWebAppOrigin:     cfg.CalendarWebAppOrigin,
		CalendarInputHandlerPath: cfg.CalendarInputHandlerPath,
	})
	return &Bot{
		log:          log,
		cfg:          cfg,
		clinic:       clinic,
		clinicDialog: clinicDialog,
		httpService: adapters.NewHttpService(
			component_name,
			&http.Server{
				Addr:    cfg.WebHandlerAddress,
				Handler: adapters.Logging(log, mux),
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
	telegram.UseRouter(ctx, &b.wg, b.log, b.bot, b.clinic, b.clinicDialog)
	go b.bot.Start()
	return nil
}

func (b *Bot) Stop(ctx context.Context) error {
	b.wg.Wait()
	b.bot.Stop()
	return b.httpService.Stop(ctx)
}
