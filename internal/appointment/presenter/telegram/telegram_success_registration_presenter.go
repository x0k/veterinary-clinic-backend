package appointment_telegram_presenter

import (
	adapters_telegram "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"gopkg.in/telebot.v3"
)

const successRegistrationText = "Вы успешно зарегистрированы!"

type SuccessRegistrationPresenter struct{}

func NewSuccessRegistrationPresenter() *SuccessRegistrationPresenter {
	return &SuccessRegistrationPresenter{}
}

func (p *SuccessRegistrationPresenter) RenderSuccessRegistration() (adapters_telegram.TextResponse, error) {
	return adapters_telegram.TextResponse{
		Text: successRegistrationText,
		Options: &telebot.SendOptions{
			ReplyMarkup: &telebot.ReplyMarkup{RemoveKeyboard: true},
		},
	}, nil
}
