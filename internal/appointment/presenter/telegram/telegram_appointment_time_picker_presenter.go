package appointment_telegram_presenter

import (
	"fmt"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/telegram"
	"gopkg.in/telebot.v3"
)

type TimePickerPresenter struct {
	stateSaver adapters.StateSaver[appointment_telegram_adapters.AppointmentSate]
}

func NewTimePickerPresenter(stateSaver adapters.StateSaver[appointment_telegram_adapters.AppointmentSate]) *TimePickerPresenter {
	return &TimePickerPresenter{
		stateSaver: stateSaver,
	}
}

func (p *TimePickerPresenter) RenderTimePicker(
	serviceId appointment.ServiceId,
	appointmentDate time.Time,
	slots appointment.SampledFreeTimeSlots,
) (telegram_adapters.TextResponses, error) {
	buttons := make([][]telebot.InlineButton, 0, len(slots)+1)
	for _, slot := range slots {
		buttons = append(buttons, []telebot.InlineButton{{
			Text:   fmt.Sprintf("%s - %s", slot.Start.String(), slot.End.String()),
			Unique: appointment_telegram_adapters.MakeAppointmentTime,
			Data: string(p.stateSaver.Save(appointment_telegram_adapters.AppointmentSate{
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
		*appointment_telegram_adapters.CancelMakeAppointmentTimeBtn.With(string(
			p.stateSaver.Save(appointment_telegram_adapters.AppointmentSate{
				ServiceId: serviceId,
				Date:      appointmentDate,
			}),
		)),
	})
	return telegram_adapters.TextResponses{{
		Text: "Выберите время:",
		Options: &telebot.SendOptions{
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: buttons,
			},
		},
	}}, nil
}
