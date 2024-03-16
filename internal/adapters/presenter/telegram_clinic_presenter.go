package presenter

import (
	"strings"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"gopkg.in/telebot.v3"
)

type TelegramClinicPresenter struct{}

func NewTelegramClinic() *TelegramClinicPresenter {
	return &TelegramClinicPresenter{}
}

func (p *TelegramClinicPresenter) RenderServices(services []entity.Service) (adapters.TelegramResponse, error) {
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
