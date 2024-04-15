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

const appointmentEventsControllerName = "appointment_pubsub_controller.AppointmentEventsController"

func NewAppointmentEvents[R any](
	log *logger.Logger,
	subs pubsub.SubscriptionsManager[appointment.EventType],
	sendAdminNotificationUseCase *appointment_use_case.SendAdminNotificationUseCase[R],
	preStopper module.PreStopper,
) module.Service {
	return module.NewService(
		appointmentEventsControllerName,
		func(ctx context.Context) error {
			l := log.With(sl.Component(appointmentEventsControllerName))
			h := func(ctx context.Context, err error) {
				if err != nil {
					l.Error(ctx, "failed to handle event", sl.Err(err))
				}
			}
			appointmentCreated := Subscribe[appointment.CreatedEvent](subs, preStopper)
			appointmentCanceled := Subscribe[appointment.CanceledEvent](subs, preStopper)
			for {
				select {
				case <-ctx.Done():
					return nil
				case e := <-appointmentCreated:
					h(ctx, sendAdminNotificationUseCase.SendAdminNotification(ctx, e))
				case e := <-appointmentCanceled:
					h(ctx, sendAdminNotificationUseCase.SendAdminNotification(ctx, e))
				}
			}
		},
	)
}
