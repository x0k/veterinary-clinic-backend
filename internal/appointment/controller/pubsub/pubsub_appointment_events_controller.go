package appointment_pubsub_controller

import (
	"context"
	"fmt"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/pubsub"
)

func NewAppointmentEvents(
	subs pubsub.SubscriptionsManager[appointment.EventType],
	preStopper module.PreStopper,
) module.Service {
	return module.NewService(
		"appointment_pubsub_controller.NewAppointmentEvents",
		func(ctx context.Context) error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case e := <-Subscribe[appointment.AppointmentCreatedEvent](subs, preStopper):
					fmt.Println(e)
				case e := <-Subscribe[appointment.AppointmentCanceledEvent](subs, preStopper):
					fmt.Println(e)
				}
			}
		},
	)
}
