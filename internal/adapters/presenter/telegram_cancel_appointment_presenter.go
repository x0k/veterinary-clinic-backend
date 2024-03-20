package presenter

import (
	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"gopkg.in/telebot.v3"
)

type TelegramCancelAppointmentPresenter struct{}

func NewTelegramCancelAppointmentPresenter() *TelegramCancelAppointmentPresenter {
	return &TelegramCancelAppointmentPresenter{}
}

func (p *TelegramCancelAppointmentPresenter) RenderCancel() (adapters.TelegramCallbackResponse, error) {
	return adapters.TelegramCallbackResponse{
		Response: &telebot.CallbackResponse{
			Text: "Ваша запись отменена",
		},
	}, nil
}
