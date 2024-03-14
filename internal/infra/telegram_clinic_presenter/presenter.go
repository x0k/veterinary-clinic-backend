package telegram_clinic_presenter

import (
	"strings"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	"gopkg.in/telebot.v3"
)

type TelegramClinicPresenter struct{}

func New() *TelegramClinicPresenter {
	return &TelegramClinicPresenter{}
}

func (p *TelegramClinicPresenter) RenderServices(services []entity.Service) (shared.TelegramResponse, error) {
	sb := strings.Builder{}
	sb.WriteString("Услуги: \n\n")
	for _, service := range services {
		sb.WriteString("*")
		sb.WriteString(service.Title)
		sb.WriteString("*\n")
		if service.Description != "" {
			sb.WriteString(shared.EscapeTelegramMarkdownString(service.Description))
			sb.WriteString("\n")
		}
		sb.WriteString(shared.EscapeTelegramMarkdownString(service.CostDescription))
		sb.WriteString("\n\n")
	}
	return shared.TelegramTextResponse{
		Text:    sb.String(),
		Options: &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2},
	}, nil
}
