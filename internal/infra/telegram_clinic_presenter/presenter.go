package telegram_clinic_presenter

import (
	"regexp"
	"strings"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	"gopkg.in/telebot.v3"
)

var escapeRegExp = regexp.MustCompile(`(\.|\-)`)

type TelegramClinicPresenter struct{}

func New() *TelegramClinicPresenter {
	return &TelegramClinicPresenter{}
}

func (p *TelegramClinicPresenter) escape(text string) string {
	return escapeRegExp.ReplaceAllString(text, "\\$1")
}

func (p *TelegramClinicPresenter) RenderServices(services []entity.Service) (shared.TelegramResponse, error) {
	sb := strings.Builder{}
	sb.WriteString("Услуги: \n\n")
	for _, service := range services {
		sb.WriteString("*")
		sb.WriteString(service.Title)
		sb.WriteString("*\n")
		if service.Description != "" {
			sb.WriteString(p.escape(service.Description))
			sb.WriteString("\n")
		}
		sb.WriteString(p.escape(service.CostDescription))
		sb.WriteString("\n\n")
	}
	return shared.TelegramResponse{
		Text:    sb.String(),
		Options: &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2},
	}, nil
}
