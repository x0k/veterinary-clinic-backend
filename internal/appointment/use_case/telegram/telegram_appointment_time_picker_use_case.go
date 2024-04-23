package appointment_telegram_use_case

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

const appointmentTimePickerUseCaseName = "appointment_telegram_use_case.AppointmentTimePickerUseCase"

type AppointmentTimePickerUseCase[R any] struct {
	log                 *logger.Logger
	schedulingService   *appointment.SchedulingService
	serviceLoader       appointment.ServiceLoader
	timePickerPresenter appointment.TimePickerPresenter[R]
	errorPresenter      appointment.ErrorPresenter[R]
}

func NewAppointmentTimePickerUseCase[R any](
	log *logger.Logger,
	schedulingService *appointment.SchedulingService,
	serviceLoader appointment.ServiceLoader,
	timePickerPresenter appointment.TimePickerPresenter[R],
	errorPresenter appointment.ErrorPresenter[R],
) *AppointmentTimePickerUseCase[R] {
	return &AppointmentTimePickerUseCase[R]{
		log:                 log.With(sl.Component(appointmentTimePickerUseCaseName)),
		schedulingService:   schedulingService,
		serviceLoader:       serviceLoader,
		timePickerPresenter: timePickerPresenter,
		errorPresenter:      errorPresenter,
	}
}

func (u *AppointmentTimePickerUseCase[R]) TimePicker(
	ctx context.Context,
	serviceId appointment.ServiceId,
	now time.Time,
	appointmentDate time.Time,
) (R, error) {
	service, err := u.serviceLoader(ctx, serviceId)
	if err != nil {
		u.log.Error(ctx, "failed to load service", sl.Err(err))
		return u.errorPresenter(err)
	}
	sampledFreeTimeSlots, err := u.schedulingService.SampledFreeTimeSlots(
		ctx,
		now,
		appointmentDate,
		service.DurationInMinutes,
	)
	if err != nil {
		u.log.Error(ctx, "failed to get sampled free time slots", sl.Err(err))
		return u.errorPresenter(err)
	}
	return u.timePickerPresenter(serviceId, appointmentDate, sampledFreeTimeSlots)
}
