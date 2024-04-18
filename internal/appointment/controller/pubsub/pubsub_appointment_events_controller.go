package appointment_pubsub_controller

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/pubsub"
)

const appointmentEventsControllerName = "appointment_pubsub_controller.AppointmentEventsController"

func NewAppointmentEvents[R any](
	subs pubsub.SubscriptionsManager[appointment.EventType],
	sendAdminNotificationUseCase *appointment_use_case.SendAdminNotificationUseCase[R],
	sendCustomerNotificationUseCase *appointment_use_case.SendCustomerNotificationUseCase[R],
	updateAppointmentsUseCase *appointment_use_case.UpdateAppointmentsStateUseCase,
	preStopper module.PreStopper,
) module.Service {
	return module.NewService(
		appointmentEventsControllerName,
		func(ctx context.Context) error {
			appointmentCreated := Subscribe[appointment.CreatedEvent](subs, preStopper)
			appointmentCanceled := Subscribe[appointment.CanceledEvent](subs, preStopper)
			appointmentChanged := Subscribe[appointment.ChangedEvent](subs, preStopper)
			for {
				select {
				case <-ctx.Done():
					return nil
				case e := <-appointmentCreated:
					updateAppointmentsUseCase.AddAppointment(ctx, e.Record)
					sendAdminNotificationUseCase.SendAdminNotification(ctx, e)
				case e := <-appointmentCanceled:
					updateAppointmentsUseCase.RemoveAppointment(ctx, e.Record)
					sendAdminNotificationUseCase.SendAdminNotification(ctx, e)
				case e := <-appointmentChanged:
					sendCustomerNotificationUseCase.SendCustomerNotification(ctx, e)
				}
			}
		},
	)
}
