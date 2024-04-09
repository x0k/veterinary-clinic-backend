package appointment_telegram_presenter

import (
	adapters_telegram "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"gopkg.in/telebot.v3"
)

type ErrorPresenter struct{}

func NewErrorPresenter() *ErrorPresenter {
	return &ErrorPresenter{}
}

func (p *ErrorPresenter) RenderError(err error) (adapters_telegram.TextResponse, error) {
	// TODO: Handle domain errors
	return adapters_telegram.TextResponse{
		Text: "Что-то пошло не так.",
		Options: &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdownV2,
		},
	}, nil
}
