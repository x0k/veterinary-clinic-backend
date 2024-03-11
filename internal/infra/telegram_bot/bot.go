package telegram_bot

import (
	"context"
	"fmt"

	"github.com/x0k/veterinary-clinic-backend/internal/config"
	"github.com/x0k/veterinary-clinic-backend/internal/controller/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"gopkg.in/telebot.v3"
)

type Bot struct {
	cfg *config.Config
	bot *telebot.Bot
}

func New(cfg *config.Config, log *logger.Logger) *Bot {
	return &Bot{
		cfg: cfg,
	}
}

func (b *Bot) Name() string {
	return "telegram_bot"
}

func (b *Bot) Start(ctx context.Context) error {
	if bot, err := telebot.NewBot(telebot.Settings{
		Token: b.cfg.TelegramToken,
	}); err != nil {
		return fmt.Errorf("starting telebot: %w", err)
	} else {
		b.bot = bot
	}
	telegram.UseRouter(b.bot)
	go b.bot.Start()
	return nil
}

func (b *Bot) Stop(ctx context.Context) error {
	b.bot.Stop()
	return nil
}
