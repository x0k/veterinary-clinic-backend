package appointment_telegram_presenter

import (
	adapters_telegram "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"gopkg.in/telebot.v3"
)

type RegistrationPresenter struct {
}

func NewRegistrationPresenter() *RegistrationPresenter {
	return &RegistrationPresenter{}
}

func (p *RegistrationPresenter) RenderRegistration() (adapters_telegram.TextResponse, error) {
	return adapters_telegram.TextResponse{
		Text: "Для записи на прием, необходимо уточнить ваш номер телефона.",
		Options: &telebot.SendOptions{
			ReplyMarkup: &telebot.ReplyMarkup{
				ReplyKeyboard: [][]telebot.ReplyButton{
					{*adapters_telegram.RegisterTelegramCustomerBtn},
				},
			},
		},
	}, nil
}
