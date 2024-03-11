package telegram_bot

import (
	"context"
	"fmt"

	"github.com/x0k/veterinary-clinic-backend/internal/config"
	"github.com/x0k/veterinary-clinic-backend/internal/controller/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
	"gopkg.in/telebot.v3"
)

type Bot struct {
	cfg    *config.TelegramConfig
	bot    *telebot.Bot
	clinic *usecase.ClinicUseCase[string]
}

func New(
	cfg *config.TelegramConfig,
	clinic *usecase.ClinicUseCase[string],
) *Bot {
	return &Bot{
		cfg:    cfg,
		clinic: clinic,
	}
}

func (b *Bot) Name() string {
	return "telegram_bot"
}

func (b *Bot) Start(ctx context.Context) error {
	if bot, err := telebot.NewBot(telebot.Settings{
		Token: b.cfg.Token,
		Poller: &telebot.LongPoller{
			Timeout: b.cfg.PollerTimeout,
		},
	}); err != nil {
		return fmt.Errorf("starting telebot: %w", err)
	} else {
		b.bot = bot
	}
	telegram.UseRouter(ctx, b.bot, b.clinic)
	go b.bot.Start()
	return nil
}

func (b *Bot) Stop(ctx context.Context) error {
	b.bot.Stop()
	return nil
}
