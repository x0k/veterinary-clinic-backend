package appointment_telegram_adapters

import (
	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"gopkg.in/telebot.v3"
)

type ErrorSender interface {
	Send(c telebot.Context, err error) error
}

type errorSender struct {
	errorPresenter appointment.ErrorPresenter[telegram_adapters.TextResponses]
}

func (s *errorSender) Send(c telebot.Context, err error) error {
	res, err := s.errorPresenter(err)
	if err != nil {
		return err
	}
	return res.Send(c)
}

func NewErrorSender(errorPresenter appointment.ErrorPresenter[telegram_adapters.TextResponses]) ErrorSender {
	return &errorSender{errorPresenter: errorPresenter}
}
