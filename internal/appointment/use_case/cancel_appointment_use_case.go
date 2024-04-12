package appointment_use_case

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

const cancelAppointmentUseCaseName = "appointment_use_case.CancelAppointmentUseCase"

type CancelAppointmentUseCase[R any] struct {
	log                        *logger.Logger
	schedulingService          *appointment.SchedulingService
	customerLoader             appointment.CustomerLoader
	appointmentCancelPresenter appointment.AppointmentCancelPresenter[R]
	errorPresenter             appointment.ErrorPresenter[R]
}

func NewCancelAppointmentUseCase[R any](
	log *logger.Logger,
	schedulingService *appointment.SchedulingService,
	customerLoader appointment.CustomerLoader,
	appointmentCancelPresenter appointment.AppointmentCancelPresenter[R],
	errorPresenter appointment.ErrorPresenter[R],
) *CancelAppointmentUseCase[R] {
	return &CancelAppointmentUseCase[R]{
		log:                        log.With(sl.Component(cancelAppointmentUseCaseName)),
		schedulingService:          schedulingService,
		customerLoader:             customerLoader,
		appointmentCancelPresenter: appointmentCancelPresenter,
		errorPresenter:             errorPresenter,
	}
}

func (s *CancelAppointmentUseCase[R]) CancelAppointment(
	ctx context.Context,
	customerIdentity appointment.CustomerIdentity,
) (R, error) {
	customer, err := s.customerLoader.Customer(ctx, customerIdentity)
	if err != nil {
		s.log.Error(ctx, "failed to load customer", sl.Err(err))
		return s.errorPresenter.RenderError(err)
	}
	if err = s.schedulingService.CancelAppointmentForCustomer(ctx, customer); err != nil {
		s.log.Error(ctx, "failed to cancel appointment", sl.Err(err))
		return s.errorPresenter.RenderError(err)
	}
	return s.appointmentCancelPresenter.RenderCancel()
}
