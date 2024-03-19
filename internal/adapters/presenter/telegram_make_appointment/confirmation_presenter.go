package telegram_make_appointment

import (
	"strings"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
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
	service entity.Service,
	appointmentDateTime time.Time,
) (adapters.TelegramTextResponse, error) {
	sb := strings.Builder{}
	sb.WriteString("Подтвердите запись:\n\n")
	sb.WriteString(adapters.EscapeTelegramMarkdownString(service.Title))
	sb.WriteString("\n\n")
	if service.Description != "" {
		sb.WriteString(adapters.EscapeTelegramMarkdownString(service.Description))
		sb.WriteString("\n\n")
	}
	sb.WriteString(adapters.EscapeTelegramMarkdownString(service.CostDescription))
	sb.WriteString("\n\n")
	sb.WriteString(adapters.EscapeTelegramMarkdownString(appointmentDateTime.Format("02.01.2006 15:04")))
	sb.WriteString(" \\- ")
	sb.WriteString(adapters.EscapeTelegramMarkdownString(
		appointmentDateTime.Add(time.Duration(service.DurationInMinutes) * time.Minute).Format("15:04"),
	))

	return adapters.TelegramTextResponse{
		Text: sb.String(),
		Options: &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdownV2,
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{
						*adapters.ConfirmMakeAppointmentBtn.With(string(p.stateSaver.Save(
							adapters.TelegramDatePickerState{
								ServiceId: service.Id,
								Date:      appointmentDateTime,
							},
						))),
					},
				},
			},
		},
	}, nil
}
