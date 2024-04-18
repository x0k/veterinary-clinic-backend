package telegram_adapters

import (
	"context"

	"gopkg.in/telebot.v3"
)

type Sender struct {
	bot *telebot.Bot
}

func NewSender(bot *telebot.Bot) *Sender {
	return &Sender{
		bot: bot,
	}
}

func (s *Sender) Send(_ context.Context, msg Message) error {
	return msg.Send(s.bot)
}
