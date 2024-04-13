package telegram_make_appointment

import (
	"fmt"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	"gopkg.in/telebot.v3"
)

type TimeSlotsPickerPresenter struct {
	stateSaver adapters.StateSaver[adapters.TelegramDatePickerState]
}

func NewTelegramTimeSlotsPickerPresenter(
	stateSaver adapters.StateSaver[adapters.TelegramDatePickerState],
) *TimeSlotsPickerPresenter {
	return &TimeSlotsPickerPresenter{
		stateSaver: stateSaver,
	}
}

func (p *TimeSlotsPickerPresenter) RenderTimePicker(
	serviceId shared.ServiceId,
	appointmentDate time.Time,
	slots shared.SampledFreeTimeSlots,
) (adapters.TelegramTextResponse, error) {
	buttons := make([][]telebot.InlineButton, 0, len(slots)+1)
	for _, slot := range slots {
		buttons = append(buttons, []telebot.InlineButton{{
			Text:   fmt.Sprintf("%s - %s", slot.Start.String(), slot.End.String()),
			Unique: adapters.MakeAppointmentTime,
			Data: string(p.stateSaver.Save(adapters.TelegramDatePickerState{
				ServiceId: serviceId,
				Date: time.Date(
					appointmentDate.Year(),
					appointmentDate.Month(),
					appointmentDate.Day(),
					slot.Start.Hours,
					slot.Start.Minutes,
					0,
					0,
					appointmentDate.Location(),
				),
			})),
		}})
	}
	buttons = append(buttons, []telebot.InlineButton{
		*adapters.CancelMakeAppointmentTimeBtn.With(string(
			p.stateSaver.Save(adapters.TelegramDatePickerState{
				ServiceId: serviceId,
				Date:      appointmentDate,
			}),
		)),
	})
	return adapters.TelegramTextResponse{
		Text: "Выберите время:",
		Options: &telebot.SendOptions{
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: buttons,
			},
		},
	}, nil
}
