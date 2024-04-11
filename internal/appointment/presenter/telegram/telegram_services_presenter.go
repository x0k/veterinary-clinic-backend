package appointment_telegram_presenter

import (
	"strings"

	adapters_telegram "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"gopkg.in/telebot.v3"
)

type ServicesPresenter struct{}

func NewServices() *ServicesPresenter {
	return &ServicesPresenter{}
}

func (s *ServicesPresenter) RenderServices(services []appointment.ServiceEntity) (adapters_telegram.TextResponses, error) {
	sb := strings.Builder{}
	writeServices(&sb, services)
	return adapters_telegram.TextResponses{{
		Text:    sb.String(),
		Options: &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2},
	}}, nil
}
