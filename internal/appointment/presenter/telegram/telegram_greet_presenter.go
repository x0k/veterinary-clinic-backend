package appointment_telegram_presenter

import (
	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	appointment_telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/telegram"
	"gopkg.in/telebot.v3"
)

func RenderGreeting() (telegram_adapters.TextResponses, error) {
	return telegram_adapters.TextResponses{{
		Text: telegram_adapters.EscapeMarkdownString("Привет!"),
		Options: &telebot.SendOptions{
			ParseMode:   telebot.ModeMarkdownV2,
			ReplyMarkup: appointment_telegram_adapters.BotMenu,
		},
	}}, nil
}
