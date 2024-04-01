package appointment_notion

import (
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/containers"
)

type Appointment struct {
	client        *notionapi.Client
	servicesCache *containers.Expiable[[]appointment.Service]
}

func NewAppointment(client *notionapi.Client) *Appointment {
	return &Appointment{
		client:        client,
		servicesCache: containers.NewExpiable[[]appointment.Service](time.Hour),
	}
}
