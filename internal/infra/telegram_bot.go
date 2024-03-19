package infra

import (
	"context"

	"gopkg.in/telebot.v3"
)

type TelegramBot struct {
	bot   *telebot.Bot
	start func(ctx context.Context, bot *telebot.Bot) error
}

func NewTelegramBot(
	bot *telebot.Bot,
	start func(ctx context.Context, bot *telebot.Bot) error,
) *TelegramBot {
	return &TelegramBot{
		bot:   bot,
		start: start,
	}
}

func (b *TelegramBot) Start(ctx context.Context) error {
	return b.start(ctx, b.bot)
}
