package appointment_js_use_case

import (
	"context"
	"errors"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

const activeAppointmentUseCaseName = "appointment_js_use_case.ActiveAppointmentUseCase"

type ActiveAppointmentUseCase[R any] struct {
	log                             *logger.Logger
	customerLoader                  appointment.CustomerByIdentityLoader
	customerActiveAppointmentLoader appointment.CustomerActiveAppointmentLoader
	serviceLoader                   appointment.ServiceLoader
	appointmentPresenter            appointment.AppointmentInfoPresenter[R]
	notFoundPresenter               appointment.NotFoundPresenter[R]
	errorPresenter                  appointment.ErrorPresenter[R]
}

func NewActiveAppointmentUseCase[R any](
	log *logger.Logger,
	customerLoader appointment.CustomerByIdentityLoader,
	customerActiveAppointmentLoader appointment.CustomerActiveAppointmentLoader,
	serviceLoader appointment.ServiceLoader,
	appointmentPresenter appointment.AppointmentInfoPresenter[R],
	notFoundPresenter appointment.NotFoundPresenter[R],
	errorPresenter appointment.ErrorPresenter[R],
) *ActiveAppointmentUseCase[R] {
	return &ActiveAppointmentUseCase[R]{
		log:                             log.With(sl.Component(activeAppointmentUseCaseName)),
		customerLoader:                  customerLoader,
		customerActiveAppointmentLoader: customerActiveAppointmentLoader,
		serviceLoader:                   serviceLoader,
		appointmentPresenter:            appointmentPresenter,
		notFoundPresenter:               notFoundPresenter,
		errorPresenter:                  errorPresenter,
	}
}

func (u *ActiveAppointmentUseCase[R]) ActiveAppointment(
	ctx context.Context,
	customerIdentity appointment.CustomerIdentity,
) (R, error) {
	customer, err := u.customerLoader(ctx, customerIdentity)
	if errors.Is(err, shared.ErrNotFound) {
		return u.notFoundPresenter()
	}
	if err != nil {
		u.log.Debug(ctx, "failed to load customer", sl.Err(err))
		return u.errorPresenter(err)
	}
	existedAppointment, err := u.customerActiveAppointmentLoader(ctx, customer.Id)
	if errors.Is(err, shared.ErrNotFound) {
		return u.notFoundPresenter()
	}
	if err != nil {
		u.log.Debug(ctx, "failed to load active appointment", sl.Err(err))
		return u.errorPresenter(err)
	}
	service, err := u.serviceLoader(ctx, existedAppointment.ServiceId)
	if err != nil {
		u.log.Debug(ctx, "failed to load service", sl.Err(err))
		return u.errorPresenter(err)
	}
	return u.appointmentPresenter(existedAppointment, service)
}
