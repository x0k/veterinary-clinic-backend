package appointment_telegram_presenter

import (
	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
)

func AppointmentCreatedEventPresenter(
	created appointment.AppointmentCreatedEvent,
) (telegram_adapters.Message, error) {
	return telegram_adapters.TextResponses{{}}, nil
}
