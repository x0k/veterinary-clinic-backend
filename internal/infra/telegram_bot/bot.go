package telegram_bot

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters/controller"
	"github.com/x0k/veterinary-clinic-backend/internal/infra"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
	"gopkg.in/telebot.v3"
)

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
	httpService  *infra.HttpService
	cfg          *Config
	bot          *telebot.Bot
	clinic       *usecase.ClinicUseCase[adapters.TelegramResponse]
	clinicDialog *usecase.ClinicDialogUseCase[adapters.TelegramResponse]
	stop         context.CancelFunc
}

func New(
	log *logger.Logger,
	clinic *usecase.ClinicUseCase[adapters.TelegramResponse],
	clinicDialog *usecase.ClinicDialogUseCase[adapters.TelegramResponse],
	fataler infra.Fataler,
	initDataParser controller.TelegramInitDataParser,
	cfg *Config,
) *Bot {
	mux := http.NewServeMux()
	controller.UseHttpTelegramRouter(log, mux, clinicDialog, initDataParser, &controller.HttpTelegramConfig{
		CalendarWebAppOrigin:     cfg.CalendarWebAppOrigin,
		CalendarInputHandlerPath: cfg.CalendarInputHandlerPath,
	})
	return &Bot{
		log:          log,
		cfg:          cfg,
		clinic:       clinic,
		clinicDialog: clinicDialog,
		httpService: infra.NewHttpService(
			"telegram_bot",
			&http.Server{
				Addr:    cfg.WebHandlerAddress,
				Handler: infra.Logging(log, mux),
			},
			fataler,
		),
	}
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
	ctx, b.stop = context.WithCancel(ctx)
	controller.UseTelegramBotRouter(ctx, &b.wg, b.log, b.bot, b.clinic, b.clinicDialog)
	go b.bot.Start()
	return nil
}

func (b *Bot) Stop(ctx context.Context) error {
	b.bot.Stop()
	b.stop()
	b.wg.Wait()
	return b.httpService.Stop(ctx)
}
