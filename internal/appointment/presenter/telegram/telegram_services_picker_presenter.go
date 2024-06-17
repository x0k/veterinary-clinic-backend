package appointment_telegram_presenter

import (
	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/telegram"
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

func (p *ServicesPickerPresenter) RenderServicesList(services []appointment.ServiceEntity) (telegram_adapters.TextResponses, error) {
	buttons := make([][]telebot.InlineButton, 0, len(services))
	for _, service := range services {
		buttons = append(buttons, []telebot.InlineButton{{
			Text:   service.Title,
			Unique: appointment_telegram_adapters.MakeAppointmentService,
			Data:   p.stateSaver(service.Id).String(),
		}})
	}
	return telegram_adapters.TextResponses{{
		Text: "Выберите услугу:",
		Options: &telebot.SendOptions{
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: buttons,
			},
		},
	}}, nil
}
