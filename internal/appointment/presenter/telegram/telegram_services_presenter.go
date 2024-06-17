package appointment_telegram_presenter

import (
	"strings"

	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"gopkg.in/telebot.v3"
)

func ServicesPresenter(services []appointment.ServiceEntity) (telegram_adapters.TextResponses, error) {
	sb := strings.Builder{}
	sb.WriteString("Услуги: \n\n")
	for _, service := range services {
		sb.WriteByte('*')
		sb.WriteString(service.Title)
		sb.WriteString("*\n")
		if service.Description != "" {
			sb.WriteString(telegram_adapters.EscapeMarkdownString(service.Description))
			sb.WriteString("\n")
		}
		sb.WriteString(telegram_adapters.EscapeMarkdownString(service.CostDescription))
		sb.WriteString("\n\n")
	}
	return telegram_adapters.TextResponses{{
		Text:    sb.String(),
		Options: &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2},
	}}, nil
}
