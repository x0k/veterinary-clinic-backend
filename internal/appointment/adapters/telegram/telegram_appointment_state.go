package appointment_telegram_adapters

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
)

type AppointmentSate struct {
	ServiceId appointment.ServiceId
	Date      time.Time
}
