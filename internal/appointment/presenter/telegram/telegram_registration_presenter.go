package appointment_telegram_presenter

import (
	"strconv"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	adapters_telegram "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"gopkg.in/telebot.v3"
)

type RegistrationPresenter struct {
	stateSaver adapters.StateByKeySaver[entity.TelegramUserId]
}

func NewRegistrationPresenter(
	stateSaver adapters.StateByKeySaver[entity.TelegramUserId],
) *RegistrationPresenter {
	return &RegistrationPresenter{
		stateSaver: stateSaver,
	}
}

func (p *RegistrationPresenter) RenderRegistration(telegramUserId entity.TelegramUserId) (adapters_telegram.TextResponse, error) {
	p.stateSaver.SaveByKey(
		adapters.NewStateId(strconv.FormatInt(telegramUserId.Int(), 10)),
		telegramUserId,
	)
	return adapters_telegram.TextResponse{
		Text: "Для записи на прием, необходимо уточнить ваш номер телефона.",
		Options: &telebot.SendOptions{
			ReplyMarkup: &telebot.ReplyMarkup{
				ReplyKeyboard: [][]telebot.ReplyButton{
					{*adapters_telegram.RegisterTelegramCustomerBtn},
					{*adapters_telegram.CancelRegisterTelegramCustomerBtn},
				},
			},
		},
	}, nil
}
