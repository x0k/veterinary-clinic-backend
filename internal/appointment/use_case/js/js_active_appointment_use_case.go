package appointment_js_use_case

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_js_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

const activeAppointmentUseCaseName = "appointment_js_use_case.ActiveAppointmentUseCase"

type ActiveAppointmentUseCase[R any] struct {
	log                             *logger.Logger
	customerLoader                  appointment.CustomerByIdentityLoader
	customerActiveAppointmentLoader appointment.CustomerActiveAppointmentLoader
	serviceLoader                   appointment.ServiceLoader
	appointmentPresenter            appointment.AppointmentInfoPresenter[R]
	errorPresenter                  appointment.ErrorPresenter[R]
}

func NewActiveAppointmentUseCase[R any](
	log *logger.Logger,
	customerLoader appointment.CustomerByIdentityLoader,
	customerActiveAppointmentLoader appointment.CustomerActiveAppointmentLoader,
	serviceLoader appointment.ServiceLoader,
	appointmentPresenter appointment.AppointmentInfoPresenter[R],
	errorPresenter appointment.ErrorPresenter[R],
) *ActiveAppointmentUseCase[R] {
	return &ActiveAppointmentUseCase[R]{
		log:                             log.With(sl.Component(activeAppointmentUseCaseName)),
		customerLoader:                  customerLoader,
		customerActiveAppointmentLoader: customerActiveAppointmentLoader,
		serviceLoader:                   serviceLoader,
		appointmentPresenter:            appointmentPresenter,
		errorPresenter:                  errorPresenter,
	}
}

func (u *ActiveAppointmentUseCase[R]) ActiveAppointment(
	ctx context.Context,
	userIdentityProvider appointment_js_adapters.CustomerIdentityProvider,
	userIdentity string,
) (R, error) {
	identity, err := customerIdentity(userIdentityProvider, userIdentity)
	if err != nil {
		u.log.Debug(ctx, "failed to create customer identity", sl.Err(err))
		return u.errorPresenter(err)
	}
	customer, err := u.customerLoader(ctx, identity)
	if err != nil {
		u.log.Debug(ctx, "failed to load customer", sl.Err(err))
		return u.errorPresenter(err)
	}
	existedAppointment, err := u.customerActiveAppointmentLoader(ctx, customer.Id)
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
