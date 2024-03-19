package telegram_clinic_make_appointment

import (
	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"gopkg.in/telebot.v3"
)

type TelegramServicePickerPresenter struct {
	stateSaver adapters.StateSaver[entity.ServiceId]
}

func NewTelegramServicePickerPresenter(
	stateSaver adapters.StateSaver[entity.ServiceId],
) *TelegramServicePickerPresenter {
	return &TelegramServicePickerPresenter{
		stateSaver: stateSaver,
	}
}

func (s *TelegramServicePickerPresenter) RenderServicesList(services []entity.Service) (adapters.TelegramTextResponse, error) {
	buttons := make([][]telebot.InlineButton, 0, len(services))
	for _, service := range services {
		buttons = append(buttons, []telebot.InlineButton{{
			Text:   service.Title,
			Unique: adapters.ClinicMakeAppointmentService,
			Data:   string(s.stateSaver.Save(service.Id)),
		}})
	}
	return adapters.TelegramTextResponse{
		Text: "Выберите услугу:",
		Options: &telebot.SendOptions{
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: buttons,
			},
		},
	}, nil
}
