package appointment_telegram_presenter

import (
	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"gopkg.in/telebot.v3"
)

const successRegistrationText = "Вы успешно зарегистрированы!"

type SuccessRegistrationPresenter struct {
	servicesPickerPresenter *ServicesPickerPresenter
}

func NewSuccessRegistrationPresenter(
	servicesPickerPresenter *ServicesPickerPresenter,
) *SuccessRegistrationPresenter {
	return &SuccessRegistrationPresenter{
		servicesPickerPresenter: servicesPickerPresenter,
	}
}

func (p *SuccessRegistrationPresenter) RenderSuccessRegistration(services []appointment.ServiceEntity) (telegram_adapters.TextResponses, error) {
	picker, err := p.servicesPickerPresenter.RenderServicesList(services)
	if err != nil {
		return nil, err
	}
	return append(telegram_adapters.TextResponses{
		{
			Text: successRegistrationText,
			Options: &telebot.SendOptions{
				ReplyMarkup: &telebot.ReplyMarkup{
					RemoveKeyboard: true,
				},
			},
		},
	}, picker...), nil
}
