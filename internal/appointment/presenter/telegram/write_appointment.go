package appointment_telegram_presenter

import (
	"strings"
	"time"

	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
)

func writeAppointment(
	w *strings.Builder,
	service appointment.ServiceEntity,
	appointmentDateTime time.Time,
) {
	w.WriteString(telegram_adapters.EscapeMarkdownString(service.Title))
	w.WriteString("\n\n")
	if service.Description != "" {
		w.WriteString(telegram_adapters.EscapeMarkdownString(service.Description))
		w.WriteString("\n\n")
	}
	w.WriteString(telegram_adapters.EscapeMarkdownString(service.CostDescription))
	w.WriteString("\n\n")
	w.WriteString(telegram_adapters.EscapeMarkdownString(appointmentDateTime.Format("02.01.2006 15:04")))
	w.WriteString(" \\- ")
	w.WriteString(telegram_adapters.EscapeMarkdownString(
		appointmentDateTime.Add(time.Duration(service.DurationInMinutes) * time.Minute).Format("15:04"),
	))
}
