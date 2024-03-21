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
	sb := strings.Builder{}
	sb.WriteString("Ваша запись:\n\n")
	WriteAppointment(&sb, record.Service, entity.DateTimeToGoTime(record.DateTimePeriod.Start))
	return adapters.TelegramTextResponse{
		Text: sb.String(),
		Options: &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdownV2,
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{*adapters.CancelAppointmentBtn},
				},
			},
		},
	}, nil
}
