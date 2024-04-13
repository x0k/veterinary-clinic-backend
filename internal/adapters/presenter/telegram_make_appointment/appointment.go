package telegram_make_appointment

import (
	"strings"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

func WriteAppointment(
	w *strings.Builder,
	service shared.Service,
	appointmentDateTime time.Time,
) {
	w.WriteString(adapters.EscapeTelegramMarkdownString(service.Title))
	w.WriteString("\n\n")
	if service.Description != "" {
		w.WriteString(adapters.EscapeTelegramMarkdownString(service.Description))
		w.WriteString("\n\n")
	}
	w.WriteString(adapters.EscapeTelegramMarkdownString(service.CostDescription))
	w.WriteString("\n\n")
	w.WriteString(adapters.EscapeTelegramMarkdownString(appointmentDateTime.Format("02.01.2006 15:04")))
	w.WriteString(" \\- ")
	w.WriteString(adapters.EscapeTelegramMarkdownString(
		appointmentDateTime.Add(time.Duration(service.DurationInMinutes) * time.Minute).Format("15:04"),
	))
}
