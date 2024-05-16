package appointment_use_case

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/pubsub"
)

const cancelAppointmentUseCaseName = "appointment_use_case.CancelAppointmentUseCase"

type CancelAppointmentUseCase[R any] struct {
	log                        *logger.Logger
	schedulingService          *appointment.SchedulingService
	customerLoader             appointment.CustomerByIdentityLoader
	serviceLoader              appointment.ServiceLoader
	appointmentCancelPresenter appointment.AppointmentCancelPresenter[R]
	errorPresenter             appointment.ErrorPresenter[R]
	publisher                  pubsub.Publisher[appointment.EventType]
}

func NewCancelAppointmentUseCase[R any](
	log *logger.Logger,
	schedulingService *appointment.SchedulingService,
	customerLoader appointment.CustomerByIdentityLoader,
	serviceLoader appointment.ServiceLoader,
	appointmentCancelPresenter appointment.AppointmentCancelPresenter[R],
	errorPresenter appointment.ErrorPresenter[R],
	publisher pubsub.Publisher[appointment.EventType],
) *CancelAppointmentUseCase[R] {
	return &CancelAppointmentUseCase[R]{
		log:                        log.With(sl.Component(cancelAppointmentUseCaseName)),
		schedulingService:          schedulingService,
		serviceLoader:              serviceLoader,
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
	customer, err := s.customerLoader(ctx, customerIdentity)
	if err != nil {
		s.log.Debug(ctx, "failed to load customer", sl.Err(err))
		res, err := s.errorPresenter(err)
		return false, res, err
	}
	rec, err := s.schedulingService.CancelAppointmentForCustomer(ctx, customer.Id)
	if err != nil {
		s.log.Debug(ctx, "failed to cancel appointment", sl.Err(err))
		res, err := s.errorPresenter(err)
		return false, res, err
	}
	if service, err := s.serviceLoader(ctx, rec.ServiceId); err != nil {
		s.log.Debug(ctx, "failed to load service", sl.Err(err))
	} else if err = s.publisher.Publish(appointment.NewAppointmentCanceled(rec, customer, service)); err != nil {
		s.log.Debug(ctx, "failed to publish event", sl.Err(err))
	}
	res, err := s.appointmentCancelPresenter()
	return true, res, err
}
