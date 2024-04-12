package appointment_use_case

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
)

type MakeAppointmentUseCase[R any] struct {
	schedulingService        *appointment.SchedulingService
	customerLoader           appointment.CustomerLoader
	serviceLoader            appointment.ServiceLoader
	appointmentInfoPresenter appointment.AppointmentInfoPresenter[R]
	errorPresenter           appointment.ErrorPresenter[R]
}

func NewMakeAppointmentUseCase[R any](
	schedulingService *appointment.SchedulingService,
	customerLoader appointment.CustomerLoader,
	serviceLoader appointment.ServiceLoader,
	appointmentInfoPresenter appointment.AppointmentInfoPresenter[R],
	errorPresenter appointment.ErrorPresenter[R],
) *MakeAppointmentUseCase[R] {
	return &MakeAppointmentUseCase[R]{
		schedulingService:        schedulingService,
		customerLoader:           customerLoader,
		serviceLoader:            serviceLoader,
		appointmentInfoPresenter: appointmentInfoPresenter,
		errorPresenter:           errorPresenter,
	}
}

func (s *MakeAppointmentUseCase[R]) CreateAppointment(
	ctx context.Context,
	now time.Time,
	appointmentDate time.Time,
	customerId appointment.CustomerId,
	serviceId appointment.ServiceId,
) (R, error) {
	customer, err := s.customerLoader.Customer(ctx, customerId)
	if err != nil {
		return s.errorPresenter.RenderError(err)
	}
	service, err := s.serviceLoader.Service(ctx, serviceId)
	if err != nil {
		return s.errorPresenter.RenderError(err)
	}
	appointment, err := s.schedulingService.MakeAppointment(ctx, now, appointmentDate, customer, service)
	if err != nil {
		return s.errorPresenter.RenderError(err)
	}
	return s.appointmentInfoPresenter.RenderInfo(appointment)
}
