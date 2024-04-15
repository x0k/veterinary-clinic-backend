package appointment_telegram_presenter

import (
	"strings"

	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

func writeAppointmentSummary(
	sb *strings.Builder,
	app appointment.AppointmentAggregate,
) {
	sb.WriteString(
		telegram_adapters.EscapeMarkdownString(
			shared.DateTimeToGoTime(app.DateTimePeriod().Start).
				Format("02.01.2006 15:04"),
		),
	)
	sb.WriteByte(' ')
	sb.WriteString(
		telegram_adapters.EscapeMarkdownString(
			app.Service().Title,
		),
	)
	sb.WriteString("\n")
	sb.WriteString(
		telegram_adapters.EscapeMarkdownString(
			app.Customer().Name,
		),
	)
}
