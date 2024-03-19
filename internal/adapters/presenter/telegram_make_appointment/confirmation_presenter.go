package telegram_make_appointment

import (
	"strings"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"gopkg.in/telebot.v3"
)

type TelegramConfirmationPresenter struct {
	stateSaver adapters.ManualStateSaver[adapters.TelegramDatePickerState]
}

func NewTelegramConfirmationPresenter(
	stateSaver adapters.ManualStateSaver[adapters.TelegramDatePickerState],
) *TelegramConfirmationPresenter {
	return &TelegramConfirmationPresenter{
		stateSaver: stateSaver,
	}
}

func (p *TelegramConfirmationPresenter) RenderConfirmation(
	userId entity.UserId,
	service entity.Service,
	appointmentDateTime time.Time,
) (adapters.TelegramTextResponse, error) {
	sb := strings.Builder{}
	sb.WriteString("Подтвердите запись:\n\n")
	WriteAppointment(&sb, service, appointmentDateTime)
	p.stateSaver.SaveByKey(adapters.StateId(userId), adapters.TelegramDatePickerState{
		ServiceId: service.Id,
		Date:      appointmentDateTime,
	})
	return adapters.TelegramTextResponse{
		Text: sb.String(),
		Options: &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdownV2,
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{*adapters.ConfirmMakeAppointmentBtn},
					{*adapters.CancelConfirmationAppointmentBtn},
				},
			},
		},
	}, nil
}
