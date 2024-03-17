package presenter

import (
	"strings"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"gopkg.in/telebot.v3"
)

type TelegramClinicServicesPresenter struct{}

func NewTelegramClinicServices() *TelegramClinicServicesPresenter {
	return &TelegramClinicServicesPresenter{}
}

func (p *TelegramClinicServicesPresenter) RenderServices(services []entity.Service) (adapters.TelegramTextResponse, error) {
	sb := strings.Builder{}
	sb.WriteString("Услуги: \n\n")
	for _, service := range services {
		sb.WriteByte('*')
		sb.WriteString(service.Title)
		sb.WriteString("*\n")
		if service.Description != "" {
			sb.WriteString(adapters.EscapeTelegramMarkdownString(service.Description))
			sb.WriteString("\n")
		}
		sb.WriteString(adapters.EscapeTelegramMarkdownString(service.CostDescription))
		sb.WriteString("\n\n")
	}
	return adapters.TelegramTextResponse{
		Text:    sb.String(),
		Options: &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2},
	}, nil
}
