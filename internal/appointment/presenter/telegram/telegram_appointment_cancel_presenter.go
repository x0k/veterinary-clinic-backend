package appointment_telegram_presenter

import (
	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"gopkg.in/telebot.v3"
)

func RenderAppointmentCancel() (telegram_adapters.CallbackResponse, error) {
	return telegram_adapters.CallbackResponse{
		Response: &telebot.CallbackResponse{
			Text: "Ваша запись отменена",
		},
	}, nil
}
