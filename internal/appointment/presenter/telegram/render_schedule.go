package appointment_telegram_presenter

import (
	"strings"

	adapters_telegram "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
)

func writeSchedule(sb *strings.Builder, schedule appointment.Schedule) {
	sb.WriteString("График работы на ")
	sb.WriteString(adapters_telegram.EscapeMarkdownString(
		schedule.Date.Format("02.01.2006")),
	)
	sb.WriteString(":\n\n")
	for _, period := range schedule.Entries {
		sb.WriteByte('*')
		sb.WriteString(period.Start.String())
		sb.WriteString(" \\- ")
		sb.WriteString(period.End.String())
		sb.WriteString("*\n")
		sb.WriteString(adapters_telegram.EscapeMarkdownString(period.Title))
		sb.WriteString("\n\n")
	}
	if len(schedule.Entries) == 0 {
		sb.WriteString("Нет записей\n\n")
	}
}
