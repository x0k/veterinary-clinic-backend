package appointment_telegram_presenter

import (
	"strings"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/telegram"
	"gopkg.in/telebot.v3"
)

type ConfirmationPresenter struct {
	stateSaver adapters.StateSaver[appointment_telegram_adapters.AppointmentSate]
}

func NewConfirmationPresenter(
	stateSaver adapters.StateSaver[appointment_telegram_adapters.AppointmentSate],
) *ConfirmationPresenter {
	return &ConfirmationPresenter{
		stateSaver: stateSaver,
	}
}

func (p *ConfirmationPresenter) RenderConfirmation(
	service appointment.ServiceEntity,
	appointmentDateTime time.Time,
) (telegram_adapters.TextResponses, error) {
	sb := strings.Builder{}
	sb.WriteString("Подтвердите запись:\n\n")
	writeAppointment(&sb, service, appointmentDateTime)
	stateId := string(p.stateSaver.Save(appointment_telegram_adapters.AppointmentSate{
		ServiceId: service.Id,
		Date:      appointmentDateTime,
	}))
	return telegram_adapters.TextResponses{{
		Text: sb.String(),
		Options: &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdownV2,
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{*appointment_telegram_adapters.ConfirmMakeAppointmentBtn.With(stateId)},
					{*appointment_telegram_adapters.CancelConfirmationAppointmentBtn.With(stateId)},
				},
			},
		},
	}}, nil
}
