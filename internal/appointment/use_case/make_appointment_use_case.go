package appointment_use_case

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/pubsub"
)

const makeAppointmentUseCaseName = "appointment_use_case.MakeAppointmentUseCase"

type MakeAppointmentUseCase[R any] struct {
	log                      *logger.Logger
	schedulingService        *appointment.SchedulingService
	customerLoader           appointment.CustomerByIdentityLoader
	serviceLoader            appointment.ServiceLoader
	appointmentInfoPresenter appointment.AppointmentInfoPresenter[R]
	errorPresenter           appointment.ErrorPresenter[R]
	publisher                pubsub.Publisher[appointment.EventType]
}

func NewMakeAppointmentUseCase[R any](
	log *logger.Logger,
	schedulingService *appointment.SchedulingService,
	customerLoader appointment.CustomerByIdentityLoader,
	serviceLoader appointment.ServiceLoader,
	appointmentInfoPresenter appointment.AppointmentInfoPresenter[R],
	errorPresenter appointment.ErrorPresenter[R],
	publisher pubsub.Publisher[appointment.EventType],
) *MakeAppointmentUseCase[R] {
	return &MakeAppointmentUseCase[R]{
		log:                      log.With(sl.Component(makeAppointmentUseCaseName)),
		schedulingService:        schedulingService,
		customerLoader:           customerLoader,
		serviceLoader:            serviceLoader,
		appointmentInfoPresenter: appointmentInfoPresenter,
		errorPresenter:           errorPresenter,
		publisher:                publisher,
	}
}

func (s *MakeAppointmentUseCase[R]) CreateAppointment(
	ctx context.Context,
	now time.Time,
	appointmentDate time.Time,
	customerId appointment.CustomerIdentity,
	serviceId appointment.ServiceId,
) (R, error) {
	customer, err := s.customerLoader(ctx, customerId)
	if err != nil {
		s.log.Error(ctx, "failed to load customer", sl.Err(err))
		return s.errorPresenter(err)
	}
	service, err := s.serviceLoader.Service(ctx, serviceId)
	if err != nil {
		s.log.Error(ctx, "failed to load service", sl.Err(err))
		return s.errorPresenter(err)
	}
	app, err := s.schedulingService.MakeAppointment(ctx, now, appointmentDate, customer, service)
	if err != nil {
		s.log.Error(ctx, "failed to make appointment", sl.Err(err))
		return s.errorPresenter(err)
	}
	if err := s.publisher.Publish(appointment.NewCreated(
		app,
		customer,
		service,
	)); err != nil {
		s.log.Error(ctx, "failed to publish event", sl.Err(err))
	}
	return s.appointmentInfoPresenter.RenderInfo(app, service)
}
