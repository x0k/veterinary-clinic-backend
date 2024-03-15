package telegram_bot

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters/controller"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
	"gopkg.in/telebot.v3"
)

type Config struct {
	Token         string
	PollerTimeout time.Duration
}

type Bot struct {
	log          *logger.Logger
	wg           sync.WaitGroup
	cfg          *Config
	bot          *telebot.Bot
	clinic       *usecase.ClinicUseCase[adapters.TelegramResponse]
	clinicDialog *usecase.ClinicDialogUseCase[adapters.TelegramResponse]
}

func New(
	log *logger.Logger,
	clinic *usecase.ClinicUseCase[adapters.TelegramResponse],
	clinicDialog *usecase.ClinicDialogUseCase[adapters.TelegramResponse],
	cfg *Config,
) *Bot {
	return &Bot{
		log:          log.With(slog.String("component", "infra.telegram_bot.Bot")),
		cfg:          cfg,
		clinic:       clinic,
		clinicDialog: clinicDialog,
	}
}

func (b *Bot) Start(ctx context.Context) error {
	const op = "infra.telegram_bot.Bot.Start"

	if bot, err := telebot.NewBot(telebot.Settings{
		Token: b.cfg.Token,
		Poller: &telebot.LongPoller{
			Timeout: b.cfg.PollerTimeout,
		},
	}); err != nil {
		return fmt.Errorf("%s failed to start: %w", op, err)
	} else {
		b.bot = bot
	}
	controller.UseTelegramBotRouter(ctx, &b.wg, b.log, b.bot, b.clinic, b.clinicDialog)
	context.AfterFunc(ctx, func() {
		b.bot.Stop()
	})
	b.bot.Start()
	b.wg.Wait()
	return nil
}
