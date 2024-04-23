package appointment_telegram_presenter

import (
	"strings"

	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/telegram"
	appointment_presenter "github.com/x0k/veterinary-clinic-backend/internal/appointment/presenter"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	"gopkg.in/telebot.v3"
)

func RenderAppointmentInfo(
	app appointment.RecordEntity,
	service appointment.ServiceEntity,
) (telegram_adapters.TextResponses, error) {
	status, err := appointment_presenter.RecordState(app.Status, app.IsArchived)
	if err != nil {
		return telegram_adapters.TextResponses{}, err
	}
	sb := strings.Builder{}
	sb.WriteString("Статус: ")
	sb.WriteString(telegram_adapters.EscapeMarkdownString(status))
	sb.WriteString("\n\n")
	writeAppointment(&sb, service, shared.DateTimeToGoTime(app.DateTimePeriod.Start))
	var markup *telebot.ReplyMarkup
	if app.Status == appointment.RecordAwaits {
		markup = &telebot.ReplyMarkup{
			InlineKeyboard: [][]telebot.InlineButton{
				{*appointment_telegram_adapters.CancelAppointmentBtn},
			},
		}
	}
	return telegram_adapters.TextResponses{{
		Text: sb.String(),
		Options: &telebot.SendOptions{
			ParseMode:   telebot.ModeMarkdownV2,
			ReplyMarkup: markup,
		},
	}}, nil
}
