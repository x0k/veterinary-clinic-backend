package appointment_use_case

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_event "github.com/x0k/veterinary-clinic-backend/internal/appointment/event"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/pubsub"
)

const cancelAppointmentUseCaseName = "appointment_use_case.CancelAppointmentUseCase"

type CancelAppointmentUseCase[R any] struct {
	log                        *logger.Logger
	schedulingService          *appointment.SchedulingService
	customerLoader             appointment.CustomerLoader
	appointmentCancelPresenter appointment.AppointmentCancelPresenter[R]
	errorPresenter             appointment.ErrorPresenter[R]
	publisher                  pubsub.Publisher[appointment_event.Type]
}

func NewCancelAppointmentUseCase[R any](
	log *logger.Logger,
	schedulingService *appointment.SchedulingService,
	customerLoader appointment.CustomerLoader,
	appointmentCancelPresenter appointment.AppointmentCancelPresenter[R],
	errorPresenter appointment.ErrorPresenter[R],
	publisher pubsub.Publisher[appointment_event.Type],
) *CancelAppointmentUseCase[R] {
	return &CancelAppointmentUseCase[R]{
		log:                        log.With(sl.Component(cancelAppointmentUseCaseName)),
		schedulingService:          schedulingService,
		customerLoader:             customerLoader,
		appointmentCancelPresenter: appointmentCancelPresenter,
		errorPresenter:             errorPresenter,
		publisher:                  publisher,
	}
}

// returns (canceled, response, error)
func (s *CancelAppointmentUseCase[R]) CancelAppointment(
	ctx context.Context,
	customerIdentity appointment.CustomerIdentity,
) (bool, R, error) {
	customer, err := s.customerLoader.Customer(ctx, customerIdentity)
	if err != nil {
		s.log.Error(ctx, "failed to load customer", sl.Err(err))
		res, err := s.errorPresenter.RenderError(err)
		return false, res, err
	}
	appointment, err := s.schedulingService.CancelAppointmentForCustomer(ctx, customer)
	if err != nil {
		s.log.Error(ctx, "failed to cancel appointment", sl.Err(err))
		res, err := s.errorPresenter.RenderError(err)
		return false, res, err
	}
	if err = s.publisher.Publish(appointment_event.NewAppointmentCanceled(appointment)); err != nil {
		s.log.Error(ctx, "failed to publish event", sl.Err(err))
	}
	res, err := s.appointmentCancelPresenter.RenderCancel()
	return true, res, err
}
