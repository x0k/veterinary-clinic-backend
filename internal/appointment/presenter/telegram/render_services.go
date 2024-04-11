package appointment_telegram_presenter

import (
	"strings"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	adapters_telegram "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"gopkg.in/telebot.v3"
)

func writeServices(sb *strings.Builder, services []appointment.ServiceEntity) {
	sb.WriteString("Услуги: \n\n")
	for _, service := range services {
		sb.WriteByte('*')
		sb.WriteString(service.Title)
		sb.WriteString("*\n")
		if service.Description != "" {
			sb.WriteString(adapters_telegram.EscapeMarkdownString(service.Description))
			sb.WriteString("\n")
		}
		sb.WriteString(adapters_telegram.EscapeMarkdownString(service.CostDescription))
		sb.WriteString("\n\n")
	}
}

func makeServicesButtons(
	stateSaver adapters.StateSaver[appointment.ServiceId],
	services []appointment.ServiceEntity,
) [][]telebot.InlineButton {
	buttons := make([][]telebot.InlineButton, 0, len(services))
	for _, service := range services {
		buttons = append(buttons, []telebot.InlineButton{{
			Text:   service.Title,
			Unique: adapters.MakeAppointmentService,
			Data:   stateSaver.Save(service.Id).String(),
		}})
	}
	return buttons
}
