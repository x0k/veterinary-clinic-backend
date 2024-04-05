package appointment_telegram_presenter

import (
	"strings"

	adapters_telegram "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"gopkg.in/telebot.v3"
)

type Services struct{}

func NewServices() *Services {
	return &Services{}
}

func (s *Services) RenderServices(services []appointment.ServiceEntity) (adapters_telegram.TextResponse, error) {
	sb := strings.Builder{}
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
	return adapters_telegram.TextResponse{
		Text:    sb.String(),
		Options: &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2},
	}, nil
}
