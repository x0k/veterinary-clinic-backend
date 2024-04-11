package appointment_telegram_presenter

import (
	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	adapters_telegram "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"gopkg.in/telebot.v3"
)

type ServicesPickerPresenter struct {
	stateSaver adapters.StateSaver[appointment.ServiceId]
}

func NewServicesPickerPresenter(
	stateSaver adapters.StateSaver[appointment.ServiceId],
) *ServicesPickerPresenter {
	return &ServicesPickerPresenter{
		stateSaver: stateSaver,
	}
}

func (p *ServicesPickerPresenter) RenderServicesList(services []appointment.ServiceEntity) (adapters_telegram.TextResponses, error) {
	buttons := makeServicesButtons(p.stateSaver, services)
	return adapters_telegram.TextResponses{{
		Text: "Выберите услугу:",
		Options: &telebot.SendOptions{
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: buttons,
			},
		},
	}}, nil
}
