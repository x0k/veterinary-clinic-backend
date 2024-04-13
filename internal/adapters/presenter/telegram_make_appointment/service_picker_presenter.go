package telegram_make_appointment

import (
	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	"gopkg.in/telebot.v3"
)

type TelegramServicePickerPresenter struct {
	stateSaver adapters.StateSaver[shared.ServiceId]
}

func NewTelegramServicePickerPresenter(
	stateSaver adapters.StateSaver[shared.ServiceId],
) *TelegramServicePickerPresenter {
	return &TelegramServicePickerPresenter{
		stateSaver: stateSaver,
	}
}

func (s *TelegramServicePickerPresenter) RenderServicesList(services []shared.Service) (adapters.TelegramTextResponse, error) {
	buttons := make([][]telebot.InlineButton, 0, len(services))
	for _, service := range services {
		buttons = append(buttons, []telebot.InlineButton{{
			Text:   service.Title,
			Unique: adapters.MakeAppointmentService,
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
