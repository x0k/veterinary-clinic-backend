package appointment_pubsub_controller

import (
	"context"
	"fmt"

	appointment_pubsub_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/pubsub"
	appointment_event "github.com/x0k/veterinary-clinic-backend/internal/appointment/event"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/pubsub"
)

func NewAppointment(
	subs pubsub.SubscriptionsManager[appointment_event.Type],
	preStopper module.PreStopper,
) func(context.Context) error {
	return func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case e := <-appointment_pubsub_adapters.Subscribe[appointment_event.AppointmentCreatedEvent](subs, preStopper):
				fmt.Println(e)
			}
		}
	}
}
