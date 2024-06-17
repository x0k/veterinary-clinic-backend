package appointment_telegram_presenter

import (
	"strings"

	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

func writeAppointmentSummary(
	sb *strings.Builder,
	app appointment.RecordEntity,
	customer appointment.CustomerEntity,
	service appointment.ServiceEntity,
) {
	start := shared.DateTimeToUTCTime(app.DateTimePeriod.Start)
	end := shared.DateTimeToUTCTime(app.DateTimePeriod.End)
	sb.WriteString(
		telegram_adapters.EscapeMarkdownString(
			start.Format("02.01.2006 15:04"),
		),
	)
	sb.WriteString(" \\- ")
	sb.WriteString(
		telegram_adapters.EscapeMarkdownString(
			end.Format("15:04"),
		),
	)
	sb.WriteString("\n\n")
	sb.WriteString(
		telegram_adapters.EscapeMarkdownString(
			service.Title,
		),
	)
	sb.WriteString("\n\n")
	sb.WriteString(
		telegram_adapters.EscapeMarkdownString(
			customer.Name,
		),
	)
}
