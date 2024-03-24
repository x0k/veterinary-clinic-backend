package telegram_make_appointment

import (
	"strings"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"gopkg.in/telebot.v3"
)

type TelegramAppointmentInfoPresenter struct {
}

func NewTelegramAppointmentInfoPresenter() *TelegramAppointmentInfoPresenter {
	return &TelegramAppointmentInfoPresenter{}
}

func (p *TelegramAppointmentInfoPresenter) RenderInfo(
	record entity.Record,
) (adapters.TelegramTextResponse, error) {
	status, err := entity.RecordStatusName(record.Status)
	if err != nil {
		return adapters.TelegramTextResponse{}, err
	}
	sb := strings.Builder{}
	sb.WriteString("Статус: ")
	sb.WriteString(adapters.EscapeTelegramMarkdownString(status))
	sb.WriteString("\n\n")
	WriteAppointment(&sb, record.Service, entity.DateTimeToGoTime(record.DateTimePeriod.Start))
	var markup *telebot.ReplyMarkup
	if record.Status == entity.RecordAwaits {
		markup = &telebot.ReplyMarkup{
			InlineKeyboard: [][]telebot.InlineButton{
				{*adapters.CancelAppointmentBtn},
			},
		}
	}
	return adapters.TelegramTextResponse{
		Text: sb.String(),
		Options: &telebot.SendOptions{
			ParseMode:   telebot.ModeMarkdownV2,
			ReplyMarkup: markup,
		},
	}, nil
}
