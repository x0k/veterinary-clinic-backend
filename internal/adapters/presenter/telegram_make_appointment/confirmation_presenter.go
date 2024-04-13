package telegram_make_appointment

import (
	"strings"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	"gopkg.in/telebot.v3"
)

type TelegramConfirmationPresenter struct {
	stateSaver adapters.StateSaver[adapters.TelegramDatePickerState]
}

func NewTelegramConfirmationPresenter(
	stateSaver adapters.StateSaver[adapters.TelegramDatePickerState],
) *TelegramConfirmationPresenter {
	return &TelegramConfirmationPresenter{
		stateSaver: stateSaver,
	}
}

func (p *TelegramConfirmationPresenter) RenderConfirmation(
	service shared.Service,
	appointmentDateTime time.Time,
) (adapters.TelegramTextResponse, error) {
	sb := strings.Builder{}
	sb.WriteString("Подтвердите запись:\n\n")
	WriteAppointment(&sb, service, appointmentDateTime)
	stateId := string(p.stateSaver.Save(adapters.TelegramDatePickerState{
		ServiceId: service.Id,
		Date:      appointmentDateTime,
	}))
	return adapters.TelegramTextResponse{
		Text: sb.String(),
		Options: &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdownV2,
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{*adapters.ConfirmMakeAppointmentBtn.With(stateId)},
					{*adapters.CancelConfirmationAppointmentBtn.With(stateId)},
				},
			},
		},
	}, nil
}
