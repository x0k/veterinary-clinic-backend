package presenter

import (
	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"gopkg.in/telebot.v3"
)

type TelegramGreetPresenter struct{}

func NewTelegramGreet() *TelegramGreetPresenter {
	return &TelegramGreetPresenter{}
}

func (p *TelegramGreetPresenter) RenderGreeting() (adapters.TelegramTextResponse, error) {
	return adapters.TelegramTextResponse{
		Text: adapters.EscapeTelegramMarkdownString("Привет!"),
		Options: &telebot.SendOptions{
			ParseMode:   telebot.ModeMarkdownV2,
			ReplyMarkup: adapters.BotMenu,
		},
	}, nil
}
