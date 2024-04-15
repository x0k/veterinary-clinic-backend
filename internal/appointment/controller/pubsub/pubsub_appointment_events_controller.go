package appointment_pubsub_controller

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/pubsub"
)

func NewAppointmentEvents[R any](
	log *logger.Logger,
	subs pubsub.SubscriptionsManager[appointment.EventType],
	sendAdminNotificationUseCase *appointment_use_case.SendAdminNotificationUseCase[R],
	preStopper module.PreStopper,
) module.Service {
	h := func(ctx context.Context, err error) {
		if err != nil {
			log.Error(ctx, "failed to handle event", sl.Err(err))
		}
	}
	return module.NewService(
		"appointment_pubsub_controller.NewAppointmentEvents",
		func(ctx context.Context) error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case e := <-Subscribe[appointment.AppointmentCreatedEvent](subs, preStopper):
					h(ctx, sendAdminNotificationUseCase.SendAdminNotification(ctx, e))
				case e := <-Subscribe[appointment.AppointmentCanceledEvent](subs, preStopper):
					h(ctx, sendAdminNotificationUseCase.SendAdminNotification(ctx, e))
				}
			}
		},
	)
}
