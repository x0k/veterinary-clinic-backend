package appointment_telegram_presenter

import (
	"strconv"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	appointment_telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	"gopkg.in/telebot.v3"
)

type RegistrationPresenter struct {
	stateSaver adapters.StateByKeySaver[shared.TelegramUserId]
}

func NewRegistrationPresenter(
	stateSaver adapters.StateByKeySaver[shared.TelegramUserId],
) *RegistrationPresenter {
	return &RegistrationPresenter{
		stateSaver: stateSaver,
	}
}

func (p *RegistrationPresenter) RenderRegistration(telegramUserId shared.TelegramUserId) (telegram_adapters.TextResponses, error) {
	p.stateSaver.SaveByKey(
		adapters.NewStateId(strconv.FormatInt(telegramUserId.Int(), 10)),
		telegramUserId,
	)
	return telegram_adapters.TextResponses{{
		Text: "Для записи на прием, необходимо уточнить ваш номер телефона.",
		Options: &telebot.SendOptions{
			ReplyMarkup: &telebot.ReplyMarkup{
				OneTimeKeyboard: true,
				ReplyKeyboard: [][]telebot.ReplyButton{
					{*appointment_telegram_adapters.RegisterTelegramCustomerBtn},
					{*appointment_telegram_adapters.CancelRegisterTelegramCustomerBtn},
				},
			},
		}},
	}, nil
}
